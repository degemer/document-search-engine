package index

type Counter interface {
	Count(<-chan FilteredDocument) (<-chan CountedDocument, <-chan WordsCountDoc)
	CountOne(FilteredDocument) CountedDocument
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
			wordsCount := countWords(f.Words)
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

func (sc StandardCounter) CountOne(f FilteredDocument) CountedDocument {
	return CountedDocument{Id: f.Id, WordsCount: countWords(f.Words)}
}

func countWords(words []string) map[string]int {
	wordsCount := make(map[string]int)
	for _, s := range words {
		wordsCount[s] += 1
	}
	return wordsCount
}
