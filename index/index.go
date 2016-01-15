package index

import (
	"encoding/gob"
	"log"
	"math"
	"os"
	"path/filepath"
)

const (
	CHANNEL_SIZE      int    = 50
	INDICES_DIRECTORY string = "indices"
	TFIDF             string = "tf-idf"
)

type Index interface {
	Create()
	Load()
	Save()
	Score(string) ScoredDocument
	Get(string) []DocScore
	GetAllIds() []int
	GetSum(int) float64
	GetSumSquared(int) float64
}

type DocScore struct {
	Id    int
	Score float64
}

type ByScore []DocScore

type StandardIndex struct {
	index     map[string][]DocScore
	ids       []int
	sums	  map[int]float64
	sumsSquared map[int]float64
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
	if len(ti.ids) == 0 {
		ti.ids = loadIds(TFIDF)
	}
	return ti.ids
}

func (ti *StandardIndex) GetSum(id int) float64 {
	if len(ti.sums) == 0 {
		ti.sums = loadSums(TFIDF, "sums")
	}
	return ti.sums[id]
}

func (ti *StandardIndex) GetSumSquared(id int) float64 {
	if len(ti.sumsSquared) == 0 {
		ti.sumsSquared = loadSums(TFIDF, "sumsSquared")
	}
	return ti.sumsSquared[id]
}

func (ti *TfIdf) Create() {
	countedDocuments, wordsCountDoc := ti.counter.Count(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read())))

	ti.index, ti.idf, ti.ids, ti.sums, ti.sumsSquared = CreateTfIdf(countedDocuments, wordsCountDoc)
}

func (ti *TfIdf) Load() {
	ti.index = loadIndex(TFIDF)
}

func (ti *TfIdf) Save() {
	prepareSave(TFIDF)
	saveIndex(TFIDF, ti.index)
	if len(ti.ids) != 0 {
		saveIds(TFIDF, ti.ids)
	}
	if len(ti.idf) != 0 {
		ti.saveIdf()
	}
	if len(ti.sums) != 0 {
		saveSums(TFIDF, "sums", ti.sums)
	}
	if len(ti.sumsSquared) != 0 {
		saveSums(TFIDF, "sumsSquared", ti.sumsSquared)
	}
}

func (ti *TfIdf) saveIdf() {
	idfFilePath := filepath.Join(INDICES_DIRECTORY, TFIDF, "idf")
	idfFile, err := os.Create(idfFilePath)
	if err != nil {
		log.Fatalln("Unable to create idf file ", idfFilePath, " : ", err)
	}
	idfEncoder := gob.NewEncoder(idfFile)
	idfEncoder.Encode(ti.idf)
	idfFile.Close()
}

func (ti *TfIdf) loadIdf() {
	idfFilePath := filepath.Join(INDICES_DIRECTORY, TFIDF, "idf")
	idfFile, err := os.Open(idfFilePath)
	if err != nil {
		log.Fatalln("Unable to open idf file ", idfFilePath, " : ", err)
	}
	idfEncoder := gob.NewDecoder(idfFile)
	idfEncoder.Decode(&ti.idf)
	idfFile.Close()
}

func (ti *TfIdf) Score(doc string) ScoredDocument {
	if len(ti.idf) == 0 {
		ti.loadIdf()
	}
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

func CreateTfIdf(countedDocuments <-chan CountedDocument, wordsCountDoc <-chan WordsCountDoc) (map[string][]DocScore, IdfWords, []int, map[int]float64, map[int]float64) {
	tfDocuments := []TfDocument{}
	idfWords := Idf(wordsCountDoc)
	for tfDoc := range Tf(countedDocuments) {
		// Or put it on disk
		tfDocuments = append(tfDocuments, tfDoc)
	}
	idfWord := <-idfWords
	index := make(map[string][]DocScore)
	ids := []int{}
	sums := make(map[int]float64)
	sumsSquared := make(map[int]float64)
	for _, tfDoc := range tfDocuments {
		ids = append(ids, tfDoc.Id)
		for word, freq := range tfDoc.WordsFrequency {
			score := freq * idfWord[word]
			index[word] = append(index[word], DocScore{Id: tfDoc.Id, Score: score})
			sums[tfDoc.Id] += score
			sumsSquared[tfDoc.Id] += score * score
		}
	}
	return index, idfWord, ids, sums, sumsSquared
}

func wordsTfFrequency(wordsCount map[string]int) map[string]float64 {
	wordsFrequency := make(map[string]float64)
	for word, numberWords := range wordsCount {
		wordsFrequency[word] = 1.0 + math.Log10(float64(numberWords))
	}
	return wordsFrequency
}

func prepareSave(filePath string) {
	os.MkdirAll(filepath.Join(INDICES_DIRECTORY, filePath), 0755)
}

func saveIndex(filePath string, index map[string][]DocScore) {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, "index")
	indexFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Unable to create index file ", filePath, " : ", err)
	}
	defer indexFile.Close()
	indexEncoder := gob.NewEncoder(indexFile)
	indexEncoder.Encode(index)
}

func saveIds(filePath string, ids []int) {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, "ids")
	idsFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Unable to create ids file ", filePath, " : ", err)
	}
	defer idsFile.Close()
	idsEncoder := gob.NewEncoder(idsFile)
	idsEncoder.Encode(ids)
}

func saveSums(filePath string, fileName string, sums map[int]float64) {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, fileName)
	sumsFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalln("Unable to create ids file ", filePath, " : ", err)
	}
	defer sumsFile.Close()
	idsEncoder := gob.NewEncoder(sumsFile)
	idsEncoder.Encode(sums)
}

func loadIndex(filePath string) map[string][]DocScore {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, "index")
	index := make(map[string][]DocScore)
	indexFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Unable to open index file ", filePath, " : ", err)
	}
	defer indexFile.Close()
	indexEncoder := gob.NewDecoder(indexFile)
	indexEncoder.Decode(&index)

	return index
}

func loadIds(filePath string) (ids []int) {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, "ids")
	idsFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Unable to open index file ", filePath, " : ", err)
	}
	defer idsFile.Close()
	idsEncoder := gob.NewDecoder(idsFile)
	idsEncoder.Decode(&ids)

	return
}

func loadSums(filePath string, fileName string) map[int]float64 {
	filePath = filepath.Join(INDICES_DIRECTORY, filePath, fileName)
	sums := make(map[int]float64)
	sumsFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalln("Unable to open index file ", filePath, " : ", err)
	}
	defer sumsFile.Close()
	idsEncoder := gob.NewDecoder(sumsFile)
	idsEncoder.Decode(&sums)

	return sums
}

func (r ByScore) Len() int      { return len(r) }
func (r ByScore) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool {
	return r[i].Score > r[j].Score || r[i].Score == r[j].Score && r[i].Id < r[j].Id
}
