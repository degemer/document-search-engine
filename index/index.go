package index

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CHANNEL_SIZE      int    = 50
	INDICES_DIRECTORY string = "indices"
	TFIDF             string = "tf-idf"
	TFIDFNORM         string = "tf-idf-norm"
	TFNORM            string = "tf-norm"
)

type Index interface {
	Create()
	Load() error
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
	index         map[string][]DocScore
	ids           []int
	sums          map[int]float64
	sumsSquared   map[int]float64
	reader        Reader
	tokenizer     Tokenizer
	filter        Filter
	counter       Counter
	saveDirectory string
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
	switch name {
	case "tf-idf-norm":
		temp := new(TfIdfNorm)
		temp.reader = NewReader(options)
		temp.tokenizer = NewTokenizer(options)
		temp.filter = NewFilter(options)
		temp.counter = NewCounter(options)
		temp.saveDirectory = filepath.Join(INDICES_DIRECTORY, TFIDFNORM)
		return temp
	case "tf-norm":
		temp := new(TfNorm)
		temp.reader = NewReader(options)
		temp.tokenizer = NewTokenizer(options)
		temp.filter = NewFilter(options)
		temp.counter = NewCounter(options)
		temp.saveDirectory = filepath.Join(INDICES_DIRECTORY, TFNORM)
		return temp
	}
	temp := new(TfIdf)
	temp.reader = NewReader(options)
	temp.tokenizer = NewTokenizer(options)
	temp.filter = NewFilter(options)
	temp.counter = NewCounter(options)
	temp.saveDirectory = filepath.Join(INDICES_DIRECTORY, TFIDF)

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

func (ti *StandardIndex) GetSum(id int) float64 {
	return ti.sums[id]
}

func (ti *StandardIndex) GetSumSquared(id int) float64 {
	return ti.sumsSquared[id]
}

func prepareSave(filePath string) {
	os.MkdirAll(filePath, 0755)
}

func saveIndex(filePath string, index map[string][]DocScore) {
	filePath = filepath.Join(filePath, "index")
	indexFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Unable to create index file ", filePath, " : ", err)
		os.Exit(1)
	}
	defer indexFile.Close()
	indexEncoder := gob.NewEncoder(indexFile)
	indexEncoder.Encode(index)
}

func saveIds(filePath string, ids []int) {
	filePath = filepath.Join(filePath, "ids")
	idsFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Unable to create ids file ", filePath, " : ", err)
		os.Exit(1)
	}
	defer idsFile.Close()
	idsEncoder := gob.NewEncoder(idsFile)
	idsEncoder.Encode(ids)
}

func saveSums(filePath string, fileName string, sums map[int]float64) {
	filePath = filepath.Join(filePath, fileName)
	sumsFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Unable to create ids file ", filePath, " : ", err)
		os.Exit(1)
	}
	defer sumsFile.Close()
	idsEncoder := gob.NewEncoder(sumsFile)
	idsEncoder.Encode(sums)
}

func loadIndex(filePath string) (map[string][]DocScore, error) {
	filePath = filepath.Join(filePath, "index")
	index := make(map[string][]DocScore)
	indexFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer indexFile.Close()
	indexEncoder := gob.NewDecoder(indexFile)
	indexEncoder.Decode(&index)

	return index, nil
}

func loadIds(filePath string) (ids []int) {
	filePath = filepath.Join(filePath, "ids")
	idsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open index file ", filePath, " : ", err)
		os.Exit(1)
	}
	defer idsFile.Close()
	idsEncoder := gob.NewDecoder(idsFile)
	idsEncoder.Decode(&ids)

	return
}

func loadSums(filePath string, fileName string) map[int]float64 {
	filePath = filepath.Join(filePath, fileName)
	sums := make(map[int]float64)
	sumsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Unable to open index file ", filePath, " : ", err)
		os.Exit(1)
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
