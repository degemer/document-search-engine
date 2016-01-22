package index

import (
	"github.com/reiver/go-porterstemmer"
)

const STEMMER_WORKERS = 2

type Stemmer interface {
	Stem(<-chan FilteredDocument) <-chan StemmedDocument
	StemOne(FilteredDocument) StemmedDocument
}

type StemmedDocument struct {
	Id    int
	Words []string
}

type NoStemmer struct{}
type PorterStemmer struct{}

func NewStemmer(options map[string]string) Stemmer {
	if options["stemmer"] == "stem" {
		return PorterStemmer{}
	}
	return NoStemmer{}
}

func (ns NoStemmer) Stem(filteredDocuments <-chan FilteredDocument) <-chan StemmedDocument {
	stemmedDocuments := make(chan StemmedDocument, CHANNEL_SIZE)
	go func() {
		for r := range filteredDocuments {
			stemmedDocuments <- ns.StemOne(r)
		}
		close(stemmedDocuments)
	}()
	return stemmedDocuments
}

func (ns NoStemmer) StemOne(r FilteredDocument) StemmedDocument {
	return StemmedDocument{Id: r.Id, Words: r.Words}
}

func (ps PorterStemmer) Stem(filteredDocuments <-chan FilteredDocument) <-chan StemmedDocument {
	stemmedDocuments := make(chan StemmedDocument, CHANNEL_SIZE)
	stemmedChannel := make(chan bool)

	for i := 1; i <= STEMMER_WORKERS; i++ {
		go func() {
			for r := range filteredDocuments {
				stemmedDocuments <- ps.StemOne(r)
			}
			stemmedChannel <- true
		}()
	}
	go func() {
		for i := 1; i <= STEMMER_WORKERS; i++ {
			<-stemmedChannel
		}
		close(stemmedDocuments)
	}()
	return stemmedDocuments
}

func (ps PorterStemmer) StemOne(r FilteredDocument) StemmedDocument {
	return StemmedDocument{Id: r.Id, Words: porter(r.Words)}
}

func porter(words []string) (stemmedWords []string) {
	for _, word := range words {
		stemmedWords = append(stemmedWords, string(porterstemmer.StemWithoutLowerCasing([]rune(word))))
	}
	return
}
