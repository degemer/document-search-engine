package index

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type TfIdf struct {
	StandardIndex
	idf IdfWords
}

type TfIdfNorm struct {
	TfIdf
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
