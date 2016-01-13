package search

import (
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

func BenchmarkBooleanPrefix(b *testing.B) {
    r4 := "NOT (aaa AND bbb OR ccc) OR ddd AND (eee OR NOT fff)"
    for i := 0; i < b.N; i++ {
		booleanPrefix(r4)
    }
}
