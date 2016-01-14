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
	}
	temp := new(VectorialSearch)
	temp.Index = ind
	return temp
}
