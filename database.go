package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"unicode"
)

type Reader interface {
	Read() <-chan RawDocument
}

type Tokenizer interface {
	Tokenize(<-chan RawDocument) <-chan TokenizedDocument
}

type Filter interface {
	Filter(<-chan TokenizedDocument) <-chan FilteredDocument
}

type Counter interface {
	Count(<-chan FilteredDocument) <-chan CountedDocument
}

type RawDocument struct {
	Id      int
	Content string
}

type TokenizedDocument struct {
	Id    int
	Words []string
}

type FilteredDocument struct {
	Id    int
	Words []string
}

type CountedDocument struct {
	Id         int
	WordsCount map[string]int
}

type StandardTokenizer struct{}

type CommonWordsFilter struct {
	path string
}

type StandardCounter struct{}

func (st StandardTokenizer) Tokenize(raw_documents <-chan RawDocument) <-chan TokenizedDocument {
	tokenized_documents := make(chan TokenizedDocument)
	go func() {
		for r := range raw_documents {
			tokenized_documents <- StandardTokenize(r)
		}
		close(tokenized_documents)
	}()
	return tokenized_documents
}

func (cwf CommonWordsFilter) Filter(tokenized_documents <-chan TokenizedDocument) <-chan FilteredDocument {
	filtered_documents := make(chan FilteredDocument)
	common_words := make(map[string]struct{})

	file, err := os.Open(cwf.path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		common_words[scanner.Text()] = struct{}{}
	}
	file.Close()

	go func() {
		for t := range tokenized_documents {
			filtered_documents <- CWFilter(t, common_words)
		}
		close(filtered_documents)
	}()
	return filtered_documents
}

func (sc StandardCounter) Count(filteredDocuments <-chan FilteredDocument) <-chan CountedDocument {
	countedDocuments := make(chan CountedDocument)
	go func() {
		for f := range filteredDocuments {
			countedDocuments <- CountedDocument{Id: f.Id, WordsCount: CountWords(f.Words)}
		}
		close(countedDocuments)
	}()
	return countedDocuments
}

func StandardTokenize(r RawDocument) TokenizedDocument {
	w := strings.FieldsFunc(r.Content, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
		return true
	})
	return TokenizedDocument{Id: r.Id, Words: w}
}

func CWFilter(t TokenizedDocument, common_words map[string]struct{}) FilteredDocument {
	filtered_words := []string{}
	for _, s := range t.Words {
		s = strings.ToLower(s)
		_, ok := common_words[s]
		if !ok {
			filtered_words = append(filtered_words, s)
		}
	}
	return FilteredDocument{Id: t.Id, Words: filtered_words}
}

func CountWords(words []string) map[string]int {
	wordsCount := make(map[string]int)
	for _, s := range words {
		wordsCount[s] += 1
	}
	return wordsCount
}
