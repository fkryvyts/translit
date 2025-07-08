//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fkryvyts/translit/internal/dict"
	"github.com/fkryvyts/translit/internal/parsers/jmnedict"
)

var logger = log.New(os.Stderr, "codegen: ", 0)

func main() {
	err := compileDicts()
	if err != nil {
		logger.Fatalf("Code generation failed with error: %v", err)
	}

	logger.Println("Code generation complete")
}

var cLetters = map[rune][]string{
	'a': {"あ", "ぁ", "っ", "わ", "ゎ"},
	'i': {"い", "ぃ", "っ", "ゐ"},
	'u': {"う", "ぅ", "っ"},
	'e': {"え", "ぇ", "っ", "ゑ"},
	'o': {"お", "ぉ", "っ"},
	'k': {"か", "ゕ", "き", "く", "け", "ゖ", "こ", "っ"},
	'g': {"が", "ぎ", "ぐ", "げ", "ご", "っ"},
	's': {"さ", "し", "す", "せ", "そ", "っ"},
	'z': {"ざ", "じ", "ず", "ぜ", "ぞ", "っ"},
	'j': {"ざ", "じ", "ず", "ぜ", "ぞ", "っ"},
	't': {"た", "ち", "つ", "て", "と", "っ"},
	'd': {"だ", "ぢ", "づ", "で", "ど", "っ"},
	'c': {"ち", "っ"},
	'n': {"な", "に", "ぬ", "ね", "の", "ん"},
	'h': {"は", "ひ", "ふ", "へ", "ほ", "っ"},
	'b': {"ば", "び", "ぶ", "べ", "ぼ", "っ"},
	'f': {"ふ", "っ"},
	'p': {"ぱ", "ぴ", "ぷ", "ぺ", "ぽ", "っ"},
	'm': {"ま", "み", "む", "め", "も"},
	'y': {"や", "ゃ", "ゆ", "ゅ", "よ", "ょ"},
	'r': {"ら", "り", "る", "れ", "ろ"},
	'w': {"わ", "ゐ", "ゑ", "ゎ", "を", "っ"},
	'v': {"ゔ"},
}

func compileDicts() error {
	dicts := []struct {
		srcPaths       []string
		jMnedictPath   string
		dstPath        string
		decodeCLetters bool
	}{
		{
			srcPaths: []string{
				"./internal/resources/dicts/dist/kakasidict.utf8",
				"./internal/resources/dicts/dist/unidict_adj.utf8",
				"./internal/resources/dicts/dist/unidict_noun.utf8",
			},
			jMnedictPath:   "./internal/resources/dicts/dist/JMnedict.xml",
			dstPath:        "./internal/resources/dicts/compiled/kanwa.ja.dict",
			decodeCLetters: true,
		},
		{
			srcPaths: []string{
				"./internal/resources/dicts/dist/hepburndict.utf8",
				"./internal/resources/dicts/dist/hepburnhira.utf8",
			},
			dstPath: "./internal/resources/dicts/compiled/hepburn.ja.dict",
		},
		{
			srcPaths: []string{
				"./internal/resources/dicts/dist/itaijidict.utf8",
				"./internal/resources/dicts/dist/halfkana.utf8",
			},
			dstPath: "./internal/resources/dicts/compiled/normalize.ja.dict",
		},
	}

	for _, d := range dicts {
		builder := dict.NewDictionaryBuilder()

		for _, filePath := range d.srcPaths {
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}

			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, ";;") {
					continue
				}

				parts := strings.Split(line, " ")
				if len(parts) < 2 {
					continue
				}

				translit := parts[0]
				word := parts[1]

				var letters []string

				if d.decodeCLetters {
					letters = cLetters[rune(translit[len(translit)-1])]
				}

				if len(letters) == 0 {
					builder.AddWord(dict.Word{Word: word, Translit: translit})
				}

				for _, letter := range letters {
					builder.AddWord(dict.Word{Word: word + letter, Translit: translit[:len(translit)-1] + letter})
				}
			}

			// Check for errors during the scan
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("error reading file: %w", err)
			}

			file.Close()
		}

		if d.jMnedictPath != "" {
			jdict, err := jmnedict.Parse(d.jMnedictPath)
			if err != nil {
				return fmt.Errorf("failed to parse JMnedict %w", err)
			}

			for _, entry := range jdict.Entries {
				if entry.Kanji == "" || entry.NameType != "(place)" {
					continue
				}

				builder.AddWord(dict.Word{Word: entry.Kanji, Translit: entry.Kana})
			}
		}

		builder.Build()

		err := builder.Save(d.dstPath)
		if err != nil {
			return fmt.Errorf("failed to save compiled dictionary: %w", err)
		}
	}

	return nil
}
