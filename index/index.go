package index

import (
	"encoding/gob"
	"log"
	"math"
	"os"
)

const CHANNEL_SIZE int = 50
const TFIDFFILE string = "tf-idf.index"
const IDFFILE string = "tf-idf.idf"

type Index interface {
	Create()
	Load()
	Save()
	Score(string) ScoredDocument
	Get(string) []DocScore
	GetAllIds() []int
}

type DocScore struct {
	Id    int
	Score float64
}

type StandardIndex struct {
	index     map[string][]DocScore
	ids       []int
	filePath  string
	reader    Reader
	tokenizer Tokenizer
	filter    Filter
	counter   Counter
}

type TfIdf struct {
	StandardIndex
	idf IdfWords
}

type TfDocument struct {
	Id             int
	WordsFrequency map[string]float64
}

type ScoredDocument struct {
	Id             int
	WordsFrequency map[string]float64
}

type IdfWords map[string]float64

func New(name string, options map[string]string) Index {
	temp := new(TfIdf)
	temp.filePath = options["save_path"]
	temp.reader = NewReader(options)
	temp.tokenizer = NewTokenizer(options)
	temp.filter = NewFilter(options)
	temp.counter = NewCounter(options)
	return temp
}

func (ti *StandardIndex) Get(word string) []DocScore {
	return ti.index[word]
}

func (ti *StandardIndex) GetIndex() map[string][]DocScore {
	return ti.index
}

func (ti *StandardIndex) GetAllIds() []int {
	return ti.ids
}

func (ti *TfIdf) Create() {
	countedDocuments, wordsCountDoc := ti.counter.Count(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read())))

	ti.index, ti.idf, ti.ids = CreateTfIdf(countedDocuments, wordsCountDoc)
}

func (ti *TfIdf) Load() {
	ti.index = loadIndex(TFIDFFILE)
	idfFile, err := os.Open(IDFFILE)
	if err != nil {
		log.Fatalln("Unable to open idf file ", IDFFILE, " : ", err)
	}
	idfEncoder := gob.NewDecoder(idfFile)
	idfEncoder.Decode(&ti.idf)
	idfFile.Close()
}

func (ti *TfIdf) Save() {
	saveIndex(TFIDFFILE, ti.index)
	idfFile, err := os.Create(IDFFILE)
	if err != nil {
		log.Fatalln("Unable to create idf file ", IDFFILE, " : ", err)
	}
	idfEncoder := gob.NewEncoder(idfFile)
	idfEncoder.Encode(ti.idf)
	idfFile.Close()
}

func (ti *TfIdf) Score(doc string) ScoredDocument {
	score := make(map[string]float64)
	countedDocument := ti.counter.CountOne(ti.filter.FilterOne(ti.tokenizer.TokenizeOne(RawDocument{Id: 0, Content: doc})))
	for word, freq := range wordsTfFrequency(countedDocument.WordsCount) {
		score[word] = freq * ti.idf[word]
	}
	return ScoredDocument{Id: countedDocument.Id, WordsFrequency: score}
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

func CreateTfIdf(countedDocuments <-chan CountedDocument, wordsCountDoc <-chan WordsCountDoc) (map[string][]DocScore, IdfWords, []int) {
	tfDocuments := []TfDocument{}
	idfWords := Idf(wordsCountDoc)
	for tfDoc := range Tf(countedDocuments) {
		// Or put it on disk
		tfDocuments = append(tfDocuments, tfDoc)
	}
	idfWord := <-idfWords
	index := make(map[string][]DocScore)
	ids := []int{}
	for _, tfDoc := range tfDocuments {
		ids = append(ids, tfDoc.Id)
		for word, freq := range tfDoc.WordsFrequency {
			index[word] = append(index[word], DocScore{Id: tfDoc.Id, Score: freq * idfWord[word]})
		}
	}
	return index, idfWord, ids
}

func wordsTfFrequency(wordsCount map[string]int) map[string]float64 {
	wordsFrequency := make(map[string]float64)
	for word, numberWords := range wordsCount {
		wordsFrequency[word] = 1.0 + math.Log10(float64(numberWords))
	}
	return wordsFrequency
}

func saveIndex(filePath string, index map[string][]DocScore) {
	indexFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Unable to create index file ", filePath, " : ", err)
	}
	indexEncoder := gob.NewEncoder(indexFile)
	indexEncoder.Encode(index)
	indexFile.Close()
}

func loadIndex(filePath string) map[string][]DocScore {
	index := make(map[string][]DocScore)
	indexFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Unable to open index file ", filePath, " : ", err)
	}
	indexEncoder := gob.NewDecoder(indexFile)
	indexEncoder.Decode(&index)
	indexFile.Close()

	return index
}
