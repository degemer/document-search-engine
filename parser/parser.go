package parser

import (
	"strings"
	"unicode"
)

func StandardTokenize(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
		return true
	})
}

func CWFilter(words []string, commonWords map[string]struct{}) (filteredWords []string) {
	for _, s := range words {
		s = strings.ToLower(s)
		_, ok := commonWords[s]
		if !ok {
			filteredWords = append(filteredWords, s)
		}
	}
	return
}

func CountWords(words []string) map[string]int {
	wordsCount := make(map[string]int)
	for _, s := range words {
		wordsCount[s] += 1
	}
	return wordsCount
}
