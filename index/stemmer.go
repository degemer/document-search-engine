package index

import (
	"github.com/reiver/go-porterstemmer"
)

const STEMMER_WORKERS = 4

type Stemmer interface {
	Stem(<-chan FilteredDocument) <-chan TokenizedDocument
	StemOne(FilteredDocument) TokenizedDocument
}

type NoStemmer struct{}
type PorterStemmer struct{}

func NewStemmer(options map[string]string) Stemmer {
	if options["stemmer"] == "stem" {
		return PorterStemmer{}
	}
	return NoStemmer{}
}

func (ns NoStemmer) Stem(filteredDocuments <-chan FilteredDocument) <-chan TokenizedDocument {
	stemmedDocuments := make(chan TokenizedDocument, CHANNEL_SIZE)
	go func() {
		for r := range filteredDocuments {
			stemmedDocuments <- ns.StemOne(r)
		}
		close(stemmedDocuments)
	}()
	return stemmedDocuments
}

func (ns NoStemmer) StemOne(r FilteredDocument) TokenizedDocument {
	return TokenizedDocument{Id: r.Id, Words: r.Words}
}

func (ps PorterStemmer) Stem(filteredDocuments <-chan FilteredDocument) <-chan TokenizedDocument {
	stemmedDocuments := make(chan TokenizedDocument, CHANNEL_SIZE)
	go func() {
		for r := range filteredDocuments {
			stemmedDocuments <- ps.StemOne(r)
		}
		close(stemmedDocuments)
	}()
	return stemmedDocuments
}

func (ps PorterStemmer) StemOne(r FilteredDocument) TokenizedDocument {
	return TokenizedDocument{Id: r.Id, Words: porter(r.Words)}
}

func porter(words []string) (stemmedWords []string) {
	for _, word := range words {
		stemmedWords = append(stemmedWords, string(porterstemmer.StemWithoutLowerCasing([]rune(word))))
	}
	return
}
