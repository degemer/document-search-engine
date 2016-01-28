package search

import (
	"github.com/degemer/document-search-engine/index"
)

type Searcher interface {
	Search(string) []index.DocScore
}

type StandardSearch struct {
	Index index.Index
}

func New(name string, ind index.Index) Searcher {
	switch {
	case name == "boolean":
		temp := new(BooleanSearch)
		temp.Index = ind
		return temp
	case name == "probabilistic":
		temp := new(ProbabilisticSearch)
		temp.Index = ind
		return temp
	case name == "vectorial-dice":
		temp := new(VectorialSearchSum)
		temp.score = dice
		temp.Index = ind
		return temp
	case name == "vectorial-jaccard":
		temp := new(VectorialSearchSum)
		temp.score = jaccard
		temp.Index = ind
		return temp
	case name == "vectorial-overlap":
		temp := new(VectorialSearchSum)
		temp.score = overlap
		temp.Index = ind
		return temp
	}
	temp := new(VectorialSearch)
	temp.Index = ind
	return temp
}
