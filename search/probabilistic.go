package search

import (
	"github.com/degemer/document-search-engine/index"
	"math"
	"sort"
)

type ProbabilisticSearch struct {
	StandardSearch
}

func (ps ProbabilisticSearch) Search(request string) (result []index.DocScore) {
	pi := 0.5
	piScore := math.Log10(pi / (1.0 - pi))
	numberDocuments := float64(len(ps.Index.GetAllIds()))

	resultMap := make(map[int]float64)
	scoredReq := ps.Index.Score(request)
	for word, _ := range scoredReq.WordsFrequency {
		resultsWord := ps.Index.Get(word)
		nbMatchedDocs := float64(len(resultsWord))
		qiScore := math.Log10((numberDocuments - nbMatchedDocs)/ nbMatchedDocs)
		for _, docScore := range resultsWord {
			resultMap[docScore.Id] += piScore + qiScore
		}
	}
	for id, score := range resultMap {
		result = append(
			result,
			index.DocScore{Id: id, Score: score})
	}
	sort.Sort(index.ByScore(result))
	return
}
