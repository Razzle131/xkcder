package words

import (
	"maps"
	"slices"
	"strings"
	"unicode"

	"github.com/kljensen/snowball/english"
)

func isSeparator(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsDigit(r)
}

func Norm(phrase string) []string {
	words := strings.FieldsFunc(strings.ToLower(strings.TrimSpace(phrase)), isSeparator)
	normalized := make(map[string]struct{}, len(words))
	for _, word := range words {
		if english.IsStopWord(word) || word == "" {
			continue
		}

		normalized[english.Stem(word, false)] = struct{}{}
	}

	return slices.Collect(maps.Keys(normalized))
}
