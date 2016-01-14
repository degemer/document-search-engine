package search

import (
	"github.com/degemer/document-search-engine/index"
	"math"
	"sort"
)

type VectorialSearch struct {
	StandardSearch
}

func (vs VectorialSearch) Search(request string) (result []index.DocScore) {
	scoredReq := vs.Index.Score(request)
	sumScoreReq := 0.0
	sumScoreDocs := make(map[int]float64)
	productScoreReqDocs := make(map[int]float64)
	for word, score := range scoredReq.WordsFrequency {
		sumScoreReq += score * score
		for _, docScore := range vs.Index.Get(word) {
			sumScoreDocs[docScore.Id] += docScore.Score * docScore.Score
			productScoreReqDocs[docScore.Id] += docScore.Score * score
		}
	}
	for id, sum := range sumScoreDocs {
		result = append(
			result,
			index.DocScore{Id: id, Score: cosSim(productScoreReqDocs[id], sum, sumScoreReq)})
	}
	sort.Sort(index.ByScore(result))
	return
}

func cosSim(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (math.Sqrt(sumDoc) * math.Sqrt(sumReq))
}
