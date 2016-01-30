package search

import (
	"github.com/degemer/document-search-engine/index"
	"testing"
)

func BenchmarkProbabilisticSearch(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	i := index.New("tf-idf", options)
	if err := i.Load(); err != nil {
		i.Create()
	}
	searcher := New("probabilistic", i)
	request := `I'm interested in mechanisms for communicating between disjoint processes,
possibly, but not exclusively, in a distributed environment.  I would
rather see descriptions of complete mechanisms, with or without implementations,
as opposed to theoretical work on the abstract problem.  Remote procedure
calls and message-passing are examples of my interests.`
	b.ResetTimer()
	for ind := 0; ind < b.N; ind++ {
		searcher.Search(request)
	}
}

func BenchmarkProbabilisticSearchStem(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	i := index.New("tf-idf-stem", options)
	if err := i.Load(); err != nil {
		i.Create()
	}
	searcher := New("probabilistic", i)
	request := `I'm interested in mechanisms for communicating between disjoint processes,
possibly, but not exclusively, in a distributed environment.  I would
rather see descriptions of complete mechanisms, with or without implementations,
as opposed to theoretical work on the abstract problem.  Remote procedure
calls and message-passing are examples of my interests.`
	b.ResetTimer()
	for ind := 0; ind < b.N; ind++ {
		searcher.Search(request)
	}
}
