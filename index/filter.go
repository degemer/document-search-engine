package index

import (
	"bufio"
	"github.com/degemer/document-search-engine/parser"
	"log"
	"os"
)

type Filter interface {
	Filter(<-chan TokenizedDocument) <-chan FilteredDocument
}

type FilteredDocument struct {
	Id    int
	Words []string
}

type CommonWordsFilter struct {
	Path string
}

func (cwf CommonWordsFilter) Filter(tokenizedDocuments <-chan TokenizedDocument) <-chan FilteredDocument {
	filteredDocuments := make(chan FilteredDocument)
	common_words := make(map[string]struct{})

	file, err := os.Open(cwf.Path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		common_words[scanner.Text()] = struct{}{}
	}
	file.Close()

	go func() {
		for t := range tokenizedDocuments {
			filteredDocuments <- FilteredDocument{Id: t.Id, Words: parser.CWFilter(t.Words, common_words)}
		}
		close(filteredDocuments)
	}()
	return filteredDocuments
}
