package search

import (
	"github.com/degemer/document-search-engine/index"
	"math"
	"sort"
)

type VectorialSearch struct {
	StandardSearch
}

type VectorialSearchSum struct {
	StandardSearch
	score func(float64, float64, float64) float64
}

func (vs VectorialSearch) Search(request string) (result []index.DocScore) {
	scoredReq := vs.Index.Score(request)
	sumScoreReq := 0.0
	productScoreReqDocs := make(map[int]float64)
	for word, score := range scoredReq.WordsFrequency {
		sumScoreReq += score * score
		for _, docScore := range vs.Index.Get(word) {
			productScoreReqDocs[docScore.Id] += docScore.Score * score
		}
	}
	for id, sum := range productScoreReqDocs {
		result = append(
			result,
			index.DocScore{Id: id, Score: cosSim(sum, vs.Index.GetSumSquared(id), sumScoreReq)})
	}
	sort.Sort(index.ByScore(result))
	return
}

func (vss VectorialSearchSum) Search(request string) (result []index.DocScore) {
	scoredReq := vss.Index.Score(request)
	sumScoreReq := 0.0
	productScoreReqDocs := make(map[int]float64)
	for word, score := range scoredReq.WordsFrequency {
		sumScoreReq += score
		for _, docScore := range vss.Index.Get(word) {
			productScoreReqDocs[docScore.Id] += docScore.Score * score
		}
	}
	for id, sum := range productScoreReqDocs {
		result = append(
			result,
			index.DocScore{Id: id, Score: vss.score(sum, vss.Index.GetSum(id), sumScoreReq)})
	}
	sort.Sort(index.ByScore(result))
	return
}

func cosSim(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (math.Sqrt(sumDoc) * math.Sqrt(sumReq))
}

func dice(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (sumDoc + sumReq)
}

func jaccard(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (sumDoc + sumReq - sumProduct)
}

func overlap(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	min := sumDoc
	if sumReq < sumDoc {
		min = sumReq
	}
	return sumProduct / min
}
