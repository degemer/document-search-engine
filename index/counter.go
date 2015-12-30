package index

import "github.com/degemer/document-search-engine/parser"

type Counter interface {
	Count(<-chan FilteredDocument) <-chan CountedDocument
}

type CountedDocument struct {
	Id         int
	WordsCount map[string]int
}

type StandardCounter struct{}

func (sc StandardCounter) Count(filteredDocuments <-chan FilteredDocument) <-chan CountedDocument {
	countedDocuments := make(chan CountedDocument)
	go func() {
		for f := range filteredDocuments {
			countedDocuments <- CountedDocument{Id: f.Id, WordsCount: parser.CountWords(f.Words)}
		}
		close(countedDocuments)
	}()
	return countedDocuments
}
