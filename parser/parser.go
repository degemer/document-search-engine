package parser

import (
	"strings"
	"strconv"
	"unicode"
	"log"
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

func CacmDoc(doc string) (int, string) {
	baseValues := []string{"T", "W", "B", "A", "K", "C", "N", "X"}
	presentValues := []string{}
	indValues := []int{}
	content := ""

	for _, val := range baseValues {
		if ind := strings.Index(doc, "\n."+val); ind != -1 {
			presentValues = append(presentValues, val)
			indValues = append(indValues, ind)
		}
	}

	id, err := strconv.Atoi(doc[3:indValues[0]])
	if err != nil {
		log.Fatalln("Unable to convert id ", doc[3:indValues[0]], "to int: ", err)
	}

	for i, val := range presentValues {
		if val == "T" || val == "W" || val == "K" {
			content += doc[indValues[i]+3 : indValues[i+1]]
		}
	}
	return id, content
}
