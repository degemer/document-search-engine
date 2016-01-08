package search

import (
	"github.com/degemer/document-search-engine/index"
)

type Result struct {
	Id    int
	Score float64
}

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
	return
}
