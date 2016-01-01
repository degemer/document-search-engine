package index

import (
	"math"
	"fmt"
)

const CHANNEL_SIZE int = 50

type Index interface {
	Create()
	PrintIndex()
	// Load()
	// Save()
}

type DocFreq struct {
	Id   int
	Freq float64
}

type TfIdf struct {
	options  map[string]string
	filePath string
	index    map[string][]DocFreq
}

type TfDocument struct {
	Id             int
	WordsFrequency map[string]float64
}

type IdfWords map[string]float64

func New(name string, options map[string]string) (ind Index) {
	temp := new(TfIdf)
	temp.options = options
	temp.filePath = options["save_path"]
	return temp
}

func (index *TfIdf) Create() {
	reader := NewReader(index.options)
	tokenizer := NewTokenizer(index.options)
	filter := NewFilter(index.options)
	counter := NewCounter(index.options)
	countedDocuments, wordsCountDoc := counter.Count(filter.Filter(tokenizer.Tokenize(reader.Read())))

	index.index = CreateTfIdf(countedDocuments, wordsCountDoc)
}

func (index *TfIdf) PrintIndex() {
	for word, docFreq := range index.index {
		fmt.Println(word, docFreq)
	}
}

func Tf(countedDocuments <-chan CountedDocument) <-chan TfDocument {
	tfDocuments := make(chan TfDocument, CHANNEL_SIZE)

	go func() {
		for c := range countedDocuments {
			tfDocuments <- TfDocument{Id: c.Id, WordsFrequency: wordsTfFrequency(c.WordsCount)}
		}
		close(tfDocuments)
	}()
	return tfDocuments
}

func Idf(wordsCountDoc <-chan WordsCountDoc) <-chan IdfWords {
	idfWords := make(chan IdfWords)
	go func() {
		idfWord := IdfWords{}
		wCD := <-wordsCountDoc
		numberDocuments := float64(wCD.NumberDocuments)
		for word, occurences := range wCD.WordsOccurences {
			idfWord[word] = math.Log10(numberDocuments / float64(occurences))
		}
		idfWords <- idfWord
		close(idfWords)
	}()
	return idfWords
}

func CreateTfIdf(countedDocuments <-chan CountedDocument, wordsCountDoc <-chan WordsCountDoc) map[string][]DocFreq {
	tfDocuments := []TfDocument{}
	idfWords := Idf(wordsCountDoc)
	for tfDoc := range Tf(countedDocuments) {
		// Or put it on disk
		tfDocuments = append(tfDocuments, tfDoc)
	}
	idfWord := <-idfWords
	index := make(map[string][]DocFreq)
	for _, tfDoc := range tfDocuments {
		for word, freq := range tfDoc.WordsFrequency {
			index[word] = append(index[word], DocFreq{Id: tfDoc.Id, Freq: freq * idfWord[word]})
		}
	}
	return index
}

func wordsTfFrequency(wordsCount map[string]int) map[string]float64 {
	wordsFrequency := make(map[string]float64)
	for word, numberWords := range wordsCount {
		wordsFrequency[word] = 1.0 + math.Log10(float64(numberWords))
	}
	return wordsFrequency
}
