package index

import (
	"encoding/gob"
	"fmt"
	"math"
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
type ById []DocScore

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

type TfIdf struct {
	StandardIndex
	idf IdfWords
}

type TfIdfNorm struct {
	TfIdf
}

type TfNorm struct {
	StandardIndex
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

func prepareIndex(st StandardIndex, options map[string]string) {
	st.reader = NewReader(options)
	st.tokenizer = NewTokenizer(options)
	st.filter = NewFilter(options)
	st.counter = NewCounter(options)
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

func (ti *TfIdf) Create() {
	countedDocuments, wordsCountDoc := ti.counter.Count(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read())))

	ti.index, ti.idf, ti.ids, ti.sums, ti.sumsSquared = CreateTfIdf(countedDocuments, wordsCountDoc)
}

func (ti *TfIdfNorm) Create() {
	countedDocuments, wordsCountDoc := ti.counter.Count(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read())))

	ti.index, ti.idf, ti.ids, _, ti.sumsSquared = CreateTfIdf(countedDocuments, wordsCountDoc)
	ti.sums = make(map[int]float64)
	for word, docScores := range ti.index {
		for i, docScore := range docScores {
			ti.index[word][i] = DocScore{Id: docScore.Id, Score: docScore.Score / math.Sqrt(ti.sumsSquared[docScore.Id])}
			ti.sums[docScore.Id] += ti.index[word][i].Score
		}
	}
	ti.sumsSquared = make(map[int]float64)
}

func (ti *TfNorm) Create() {
	countedDocuments, _ := ti.counter.Count(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read())))
	tfDocuments := Tf(countedDocuments)
	ti.index = make(map[string][]DocScore)
	ti.sums = make(map[int]float64)
	ti.sumsSquared = make(map[int]float64)
	for tfDoc := range tfDocuments {
		ti.ids = append(ti.ids, tfDoc.Id)
		max := 0.0
		for _, freq := range tfDoc.WordsFrequency {
			if freq > max {
				max = freq
			}
		}
		for word, freq := range tfDoc.WordsFrequency {
			score := freq / max
			ti.index[word] = append(ti.index[word], DocScore{Id: tfDoc.Id, Score: score})
			ti.sums[tfDoc.Id] += score
			ti.sumsSquared[tfDoc.Id] += score * score
		}
	}
}

func (ti *TfIdfNorm) GetSumSquared(id int) float64 {
	return 1
}

func (ti *TfIdf) Load() (err error) {
	ti.index, err = loadIndex(ti.saveDirectory)
	if err != nil {
		return
	}
	ti.loadIdf()
	ti.sumsSquared = loadSums(ti.saveDirectory, "sumsSquared")
	ti.sums = loadSums(ti.saveDirectory, "sums")
	ti.ids = loadIds(ti.saveDirectory)
	return
}

func (ti *TfIdfNorm) Load() (err error) {
	ti.index, err = loadIndex(ti.saveDirectory)
	if err != nil {
		return
	}
	ti.loadIdf()
	ti.sums = loadSums(ti.saveDirectory, "sums")
	ti.ids = loadIds(ti.saveDirectory)
	return
}

func (ti *TfNorm) Load() (err error) {
	ti.index, err = loadIndex(ti.saveDirectory)
	if err != nil {
		return
	}
	ti.sumsSquared = loadSums(ti.saveDirectory, "sumsSquared")
	ti.sums = loadSums(ti.saveDirectory, "sums")
	ti.ids = loadIds(ti.saveDirectory)
	return
}

func (ti *TfIdf) Save() {
	prepareSave(ti.saveDirectory)
	saveIndex(ti.saveDirectory, ti.index)
	if len(ti.ids) != 0 {
		saveIds(ti.saveDirectory, ti.ids)
	}
	if len(ti.idf) != 0 {
		ti.saveIdf()
	}
	if len(ti.sums) != 0 {
		saveSums(ti.saveDirectory, "sums", ti.sums)
	}
	if len(ti.sumsSquared) != 0 {
		saveSums(ti.saveDirectory, "sumsSquared", ti.sumsSquared)
	}
}

func (ti *TfNorm) Save() {
	prepareSave(ti.saveDirectory)
	saveIndex(ti.saveDirectory, ti.index)
	if len(ti.ids) != 0 {
		saveIds(ti.saveDirectory, ti.ids)
	}
	if len(ti.sums) != 0 {
		saveSums(ti.saveDirectory, "sums", ti.sums)
	}
	if len(ti.sumsSquared) != 0 {
		saveSums(ti.saveDirectory, "sumsSquared", ti.sumsSquared)
	}
}

func (ti *TfIdf) saveIdf() {
	idfFilePath := filepath.Join(ti.saveDirectory, "idf")
	idfFile, err := os.Create(idfFilePath)
	if err != nil {
		fmt.Println("Unable to create idf file ", idfFilePath, " : ", err)
	}
	idfEncoder := gob.NewEncoder(idfFile)
	idfEncoder.Encode(ti.idf)
	idfFile.Close()
}

func (ti *TfIdf) loadIdf() {
	idfFilePath := filepath.Join(ti.saveDirectory, "idf")
	idfFile, err := os.Open(idfFilePath)
	if err != nil {
		fmt.Println("Unable to open idf file ", idfFilePath, " : ", err)
		os.Exit(1)
	}
	idfEncoder := gob.NewDecoder(idfFile)
	idfEncoder.Decode(&ti.idf)
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

func (ti *TfIdfNorm) Score(doc string) ScoredDocument {
	score := make(map[string]float64)
	countedDocument := ti.counter.CountOne(ti.filter.FilterOne(ti.tokenizer.TokenizeOne(RawDocument{Id: 0, Content: doc})))
	sumSquared := 0.0
	for word, freq := range wordsTfFrequency(countedDocument.WordsCount) {
		score[word] = freq * ti.idf[word]
		sumSquared += score[word]
	}
	for word, _ := range score {
		score[word] /= math.Sqrt(sumSquared)
	}
	return ScoredDocument{Id: countedDocument.Id, WordsFrequency: score}
}

func (ti *TfNorm) Score(doc string) ScoredDocument {
	score := make(map[string]float64)
	countedDocument := ti.counter.CountOne(ti.filter.FilterOne(ti.tokenizer.TokenizeOne(RawDocument{Id: 0, Content: doc})))
	max := 0.0
	for word, freq := range wordsTfFrequency(countedDocument.WordsCount) {
		score[word] = freq
		if freq > max {
			max = freq
		}
	}
	for word, _ := range score {
		score[word] /= max
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

func (r ById) Len() int      { return len(r) }
func (r ById) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ById) Less(i, j int) bool {
	return r[i].Id < r[j].Id
}
