package index

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Filter interface {
	Filter(<-chan TokenizedDocument) <-chan FilteredDocument
	FilterOne(TokenizedDocument) FilteredDocument
}

type FilteredDocument struct {
	Id    int
	Words []string
}

type CommonWordsFilter struct {
	path        string
	commonWords map[string]struct{}
}

func NewFilter(options map[string]string) Filter {
	cwf := CommonWordsFilter{path: options["common_words_path"]}
	cwf.commonWords = make(map[string]struct{})

	file, err := os.Open(cwf.path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cwf.commonWords[scanner.Text()] = struct{}{}
	}
	file.Close()
	return cwf
}

func (cwf CommonWordsFilter) Filter(tokenizedDocuments <-chan TokenizedDocument) <-chan FilteredDocument {
	filteredDocuments := make(chan FilteredDocument, CHANNEL_SIZE)

	go func() {
		for t := range tokenizedDocuments {
			filteredDocuments <- cwf.FilterOne(t)
		}
		close(filteredDocuments)
	}()
	return filteredDocuments
}

func (cwf CommonWordsFilter) FilterOne(t TokenizedDocument) FilteredDocument {
	return FilteredDocument{Id: t.Id, Words: cwFilter(t.Words, cwf.commonWords)}
}

func cwFilter(words []string, commonWords map[string]struct{}) (filteredWords []string) {
	for _, s := range words {
		s = strings.ToLower(s)
		_, ok := commonWords[s]
		if !ok {
			filteredWords = append(filteredWords, s)
		}
	}
	return
}
