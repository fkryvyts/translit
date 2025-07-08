package translit

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTransliterate(t *testing.T) {
	testCases := []struct {
		text     string
		expected string
		lng      Language
	}{
		{
			text:     "日本国民は",
			expected: "nihonkokumin ha",
			lng:      Japanese,
		},
		{
			text:     "諸国民との協和による成果と、わが国全土にわたつて自由のもたら",
			expected: "shokokumin tono kyouwa niyoru seika to,waga kuni zendo niwatatsute jiyuu nomotara",
			lng:      Japanese,
		},
		{
			text:     "角тест",
			expected: "kaku test",
			lng:      Japanese,
		},
	}

	for _, tC := range testCases {
		printMemUsage()

		res := Transliterate(tC.text, Hints{
			Language: tC.lng,
		})

		if diff := cmp.Diff(res, tC.expected); diff != "" {
			t.Errorf("mismatched output (-got,+want):\n %s", diff)
		}
	}
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
