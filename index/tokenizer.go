package index

import (
	"strings"
	"unicode"
)

type Tokenizer interface {
	Tokenize(<-chan RawDocument) <-chan TokenizedDocument
	TokenizeOne(RawDocument) TokenizedDocument
}

type TokenizedDocument struct {
	Id    int
	Words []string
}

type StandardTokenizer struct{}

func NewTokenizer(options map[string]string) Tokenizer {
	return StandardTokenizer{}
}

func (st StandardTokenizer) Tokenize(rawDocuments <-chan RawDocument) <-chan TokenizedDocument {
	tokenizedDocuments := make(chan TokenizedDocument, CHANNEL_SIZE)
	go func() {
		for r := range rawDocuments {
			tokenizedDocuments <- st.TokenizeOne(r)
		}
		close(tokenizedDocuments)
	}()
	return tokenizedDocuments
}

func (st StandardTokenizer) TokenizeOne(r RawDocument) TokenizedDocument {
	return TokenizedDocument{Id: r.Id, Words: standardTokenize(r.Content)}
}

func standardTokenize(content string) []string {
	return strings.FieldsFunc(content, func(r rune) bool {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return false
		}
		return true
	})
}
