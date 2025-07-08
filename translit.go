package translit

import (
	"log"

	anyascii "github.com/anyascii/go"
	"github.com/fkryvyts/translit/internal/dict"
	"github.com/fkryvyts/translit/internal/resources"
)

//go:generate go run ./generate.go

var (
	kanwaJaDict      = dict.NewDictionary(200)
	hepburnJaDict    = dict.NewDictionary(-1)
	normalizeJaDict  = dict.NewDictionary(-1)
	replacementDicts = map[Language][]replacementDict{
		Japanese: {
			{dict: normalizeJaDict},
			{dict: kanwaJaDict, wordSep: " "},
			{dict: hepburnJaDict},
		},
	}
)

type replacementDict struct {
	dict    *dict.Dictionary
	wordSep string
}

func init() {
	err := kanwaJaDict.LoadFromBytes(resources.KanwaJa)
	if err != nil {
		log.Fatalf("failed to initialize kanwa ja dict: %v", err)
	}

	err = hepburnJaDict.LoadFromBytes(resources.HepburnJa)
	if err != nil {
		log.Fatalf("failed to initialize khepburn ja dict: %v", err)
	}

	err = normalizeJaDict.LoadFromBytes(resources.NormalizeJa)
	if err != nil {
		log.Fatalf("failed to initialize normalize ja dict: %v", err)
	}

	resources.KanwaJa = nil
	resources.HepburnJa = nil
	resources.NormalizeJa = nil
}

type Language string

const (
	Japanese Language = "ja"
)

type Hints struct {
	Language Language
}

func Transliterate(text string, hints Hints) string {
	dicts := replacementDicts[hints.Language]

	for _, d := range dicts {
		text = d.dict.Search(text).Replace(d.wordSep)
	}

	return anyascii.Transliterate(text)
}
