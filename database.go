package main

import (
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


type StandardTokenizer struct {}

func (st StandardTokenizer) Tokenize(raw_documents <-chan RawDocument) <-chan TokenizedDocument {
	tokenized_documents := make(chan TokenizedDocument)
	go func() {
		for r := range(raw_documents) {
			tokenized_documents <- StandardTokenize(r)
		}
		close(tokenized_documents)
	}()
	return tokenized_documents
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
