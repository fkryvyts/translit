package translit

import (
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
		res := Transliterate(tC.text, Hints{
			Language: tC.lng,
		})

		if diff := cmp.Diff(res, tC.expected); diff != "" {
			t.Errorf("mismatched output (-got,+want):\n %s", diff)
		}
	}
}
