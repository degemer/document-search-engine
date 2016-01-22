package index

import (
	"math"
)

type TfNorm struct {
	StandardIndex
}

func (ti *TfNorm) Create() {
	countedDocuments, _ := ti.counter.Count(ti.stemmer.Stem(ti.filter.Filter(ti.tokenizer.Tokenize(ti.reader.Read()))))
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

func (ti *TfNorm) Score(doc string) ScoredDocument {
	score := make(map[string]float64)
	countedDocument := ti.counter.CountOne(ti.stemmer.StemOne(ti.filter.FilterOne(ti.tokenizer.TokenizeOne(RawDocument{Id: 0, Content: doc}))))
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

func wordsTfFrequency(wordsCount map[string]int) map[string]float64 {
	wordsFrequency := make(map[string]float64)
	for word, numberWords := range wordsCount {
		wordsFrequency[word] = 1.0 + math.Log10(float64(numberWords))
	}
	return wordsFrequency
}
