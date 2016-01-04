package index

import "github.com/degemer/document-search-engine/parser"

type Counter interface {
	Count(<-chan FilteredDocument) (<-chan CountedDocument, <-chan WordsCountDoc)
}

type CountedDocument struct {
	Id         int
	WordsCount map[string]int
}

type WordsCountDoc struct {
	NumberDocuments int
	WordsOccurences map[string]int
}

type StandardCounter struct{}

func NewCounter(options map[string]string) Counter {
	return StandardCounter{}
}

func (sc StandardCounter) Count(filteredDocuments <-chan FilteredDocument) (<-chan CountedDocument, <-chan WordsCountDoc) {
	countedDocuments := make(chan CountedDocument, CHANNEL_SIZE)
	wordsCountDoc := make(chan WordsCountDoc)
	N := 0
	wordsOccurences := map[string]int{}
	go func() {
		for f := range filteredDocuments {
			N += 1
			wordsCount := parser.CountWords(f.Words)
			countedDocuments <- CountedDocument{Id: f.Id, WordsCount: wordsCount}
			for words, _ := range wordsCount {
				wordsOccurences[words] += 1
			}
		}
		close(countedDocuments)

		wordsCountDoc <- WordsCountDoc{NumberDocuments: N, WordsOccurences: wordsOccurences}
		close(wordsCountDoc)
	}()
	return countedDocuments, wordsCountDoc
}
