package index

import (
	"bufio"
	"log"
	"os"
	"strings"
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

func NewFilter(options map[string]string) Filter {
	return CommonWordsFilter{Path: options["common_words_path"]}
}

func (cwf CommonWordsFilter) Filter(tokenizedDocuments <-chan TokenizedDocument) <-chan FilteredDocument {
	filteredDocuments := make(chan FilteredDocument, CHANNEL_SIZE)
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
			filteredDocuments <- FilteredDocument{Id: t.Id, Words: cwFilter(t.Words, common_words)}
		}
		close(filteredDocuments)
	}()
	return filteredDocuments
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
