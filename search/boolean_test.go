package search

import (
	"github.com/degemer/document-search-engine/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleanPrefix(t *testing.T) {
	r1 := "aaa OR bbb AND ccc"
	prefixedR1 := []string{"or", "aaa", "and", "bbb", "ccc"}
	assert.Equal(t, booleanPrefix(r1), prefixedR1)
	r2 := "aaa AND bbb OR NOT ccc"
	prefixedR2 := []string{"or", "and", "aaa", "bbb", "not", "ccc"}
	assert.Equal(t, booleanPrefix(r2), prefixedR2)
	r3 := "aaa AND (bbb OR ccc)"
	prefixedR3 := []string{"and", "aaa", "or", "bbb", "ccc"}
	assert.Equal(t, booleanPrefix(r3), prefixedR3)
	r4 := "NOT (aaa AND bbb OR ccc)"
	prefixedR4 := []string{"not", "or", "and", "aaa", "bbb", "ccc"}
	assert.Equal(t, booleanPrefix(r4), prefixedR4)
}

func TestIntersect(t *testing.T) {
	i1 := []index.DocScore{index.DocScore{Id: 1, Score: 5}, index.DocScore{Id: 3, Score: 4}, index.DocScore{Id: 4, Score: 2}}
	i2 := []index.DocScore{index.DocScore{Id: 2, Score: 5}, index.DocScore{Id: 3, Score: 3}}
	result := []index.DocScore{index.DocScore{Id: 3, Score: 7}}
	assert.Equal(t, Intersect(i1, i2), result)
}

func TestUnion(t *testing.T) {
	i1 := []index.DocScore{index.DocScore{Id: 1, Score: 5}, index.DocScore{Id: 3, Score: 4}, index.DocScore{Id: 4, Score: 2}}
	i2 := []index.DocScore{index.DocScore{Id: 2, Score: 5}, index.DocScore{Id: 3, Score: 3}}
	result := []index.DocScore{index.DocScore{Id: 1, Score: 5}, index.DocScore{Id: 2, Score: 5}, index.DocScore{Id: 3, Score: 7}, index.DocScore{Id: 4, Score: 2}}
	assert.Equal(t, union(i1, i2), result)
}

func TestNot(t *testing.T) {
	docScores := []index.DocScore{index.DocScore{Id: 2, Score: 5}}
	ids := []int{1, 2, 3, 4}
	result := []index.DocScore{index.DocScore{Id: 1, Score: 0}, index.DocScore{Id: 3, Score: 0}, index.DocScore{Id: 4, Score: 0}}
	assert.Equal(t, not(docScores, ids), result)
}

func BenchmarkBooleanPrefix(b *testing.B) {
	r4 := "NOT (aaa AND bbb OR ccc) OR ddd AND (eee OR NOT fff)"
	for i := 0; i < b.N; i++ {
		booleanPrefix(r4)
	}
}

func BenchmarkBooleanSearch(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	i := index.New("tf-idf", options)
	if err := i.Load(); err != nil {
		i.Create()
	}
	searcher := New("boolean", i)
	request := "NOT (computer AND analysis OR ibm) OR language"
	b.ResetTimer()
	for ind := 0; ind < b.N; ind++ {
		searcher.Search(request)
	}
}
