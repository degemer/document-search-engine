package search

import (
	"bytes"
	"github.com/degemer/document-search-engine/index"
	"math"
	"sort"
	"strings"
)

const (
	BOOL_OR  = "or"
	BOOL_AND = "and"
	BOOL_NOT = "not"
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

type BooleanSearch struct {
	StandardSearch
}

func New(name string, ind index.Index) Searcher {
	switch {
	case name == "bs":
		temp := new(BooleanSearch)
		temp.Index = ind
		return temp
	}
	temp := new(VectorialSearch)
	temp.Index = ind
	return temp
}

func (vs VectorialSearch) Search(request string) (results []Result) {
	scoredReq := vs.Index.Score(request)
	sumScoreReq := 0.0
	sumScoreDocs := make(map[int]float64)
	productScoreReqDocs := make(map[int]float64)
	for word, score := range scoredReq.WordsFrequency {
		sumScoreReq += score * score
		for _, docFreq := range vs.Index.Get(word) {
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

func (bs BooleanSearch) Search(request string) (results []Result) {
	prefixedRequest := booleanPrefix(request)
	if prefixedRequest != nil {
		return
	}
	return
}

func cosSim(sumProduct float64, sumDoc float64, sumReq float64) float64 {
	return sumProduct / (math.Sqrt(sumDoc) * math.Sqrt(sumReq))
}

func (r ByScore) Len() int      { return len(r) }
func (r ByScore) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r ByScore) Less(i, j int) bool {
	return r[i].Score > r[j].Score || r[i].Score == r[j].Score && r[i].Id < r[j].Id
}

func booleanPrefix(request string) (prefixedRequest []string) {
	// No real postfix, since it's already reversed
	postfixedRequest := []string{}
	opStack := new(Stack)
	for _, val := range prepareBooleanRequest(request) {
		switch {
		case val != BOOL_NOT && val != BOOL_AND && val != BOOL_OR && val != "(" && val != ")":
			postfixedRequest = append(postfixedRequest, val)
		case val == "(":
			opStack.Push(val)
		case val == ")":
			for opStack.Top() != "(" {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Pop()
		case val == BOOL_NOT:
			for opStack.Top() == BOOL_NOT {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		case val == BOOL_AND:
			for opStack.Top() == BOOL_NOT || opStack.Top() == BOOL_AND {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		case val == BOOL_OR:
			for opStack.Top() == BOOL_NOT || opStack.Top() == BOOL_AND || opStack.Top() == BOOL_OR {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		}
	}
	for opStack.Top() != "" {
		postfixedRequest = append(postfixedRequest, opStack.Pop())
	}
	return reverse(postfixedRequest)
}

func prepareBooleanRequest(request string) []string {
	preparedRequest := new(bytes.Buffer)
	for _, val := range strings.ToLower(request) {
		if val == '(' {
			preparedRequest.WriteString(" ) ")
		} else if val == ')' {
			preparedRequest.WriteString(" ( ")
		} else {
			preparedRequest.WriteRune(val)
		}
	}
	return reverse(strings.Fields(preparedRequest.String()))
}

func reverse(r []string) []string {
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}
