package index

import "github.com/degemer/document-search-engine/parser"

type Tokenizer interface {
	Tokenize(<-chan RawDocument) <-chan TokenizedDocument
}

type TokenizedDocument struct {
	Id    int
	Words []string
}

type StandardTokenizer struct{}

func (st StandardTokenizer) Tokenize(rawDocuments <-chan RawDocument) <-chan TokenizedDocument {
	tokenizedDocuments := make(chan TokenizedDocument)
	go func() {
		for r := range rawDocuments {
			tokenizedDocuments <- TokenizedDocument{Id: r.Id, Words: parser.StandardTokenize(r.Content)}
		}
		close(tokenizedDocuments)
	}()
	return tokenizedDocuments
}
