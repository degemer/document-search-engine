package search

import (
	"github.com/degemer/document-search-engine/index"
	"math"
	"sort"
)

type Result struct {
	Id    int
	Score float64
}

type ByScore []Result

type Searcher interface {
	Search(string) []Result
}

type StandardSearch struct {
	Index index.Index
}

type VectorialSearch struct {
	StandardSearch
}

func New(name string, ind index.Index) Searcher {
	temp := new(VectorialSearch)
	temp.Index = ind
	return temp
}

func (vs VectorialSearch) Search(request string) (results []Result) {
	scoredReq := vs.Index.Score(request)
	sumScoreReq := 0.0
	sumScoreDocs := make(map[int]float64)
	productScoreReqDocs := make(map[int]float64)
	for word, score := range(scoredReq.WordsFrequency) {
		sumScoreReq += score * score
		for _, docFreq := range(vs.Index.Get(word)) {
			sumScoreDocs[docFreq.Id] += docFreq.Freq * docFreq.Freq
			productScoreReqDocs[docFreq.Id] += docFreq.Freq * score
		}
	}
	for id, sum := range sumScoreDocs {
		results = append(
			results,
			Result{Id: id, Score: cosSim(productScoreReqDocs[id], sum, sumScoreReq)})
	}
	sort.Sort(ByScore(results))
	return
}

func cosSim(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (math.Sqrt(sumDoc) * math.Sqrt(sumReq))
}

func (r ByScore) Len() int { return len(r) }
func (r ByScore) Swap(i, j int) {r[i], r[j] = r[j], r[i]}
func (r ByScore) Less(i, j int) bool {return r[i].Score > r[j].Score || r[i].Score == r[j].Score && r[i].Id < r[j].Id}
