package search

import (
	"fmt"
	"github.com/degemer/document-search-engine/index"
	"os"
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
	if name != "" && name != "vectorial" {
		fmt.Println(
			"Index has to be one of:",
			"boolean",
			"probabilistic",
			"vectorial",
			"vectorial-dice",
			"vectorial-jaccard",
			"vectorial-overlap",
		)
		fmt.Println("Received:", name)
		os.Exit(1)
	}
	temp := new(VectorialSearch)
	temp.Index = ind
	return temp
}
