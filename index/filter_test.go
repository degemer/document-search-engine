package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCWFilter(t *testing.T) {
	words := []string{
		"An", "On", "Line", "Program", "for", "Non", "Numerical", "Algebra", "The", "goal",
		"of", "this", "program", "is", "to", "make", "a", "step", "toward", "te", "design",
		"of", "an", "automated", "mathematical", "assistant", "Some", "requirements", "for",
	}
	filters := map[string]struct{}{"on": struct{}{}, "a": struct{}{}, "algebra": struct{}{}}
	filtered_words := []string{
		"an", "line", "program", "for", "non", "numerical", "the", "goal", "of",
		"this", "program", "is", "to", "make", "step", "toward", "te",
		"design", "of", "an", "automated", "mathematical", "assistant",
		"some", "requirements", "for",
	}

	assert.Equal(t, cwFilter(words, filters), filtered_words, "Filter incorrect")
}

func BenchmarkCWFilter(b *testing.B) {
	words := []string{
		"An", "On", "Line", "Program", "for", "Non", "Numerical", "Algebra", "The", "goal",
		"of", "this", "program", "is", "to", "make", "a", "step", "toward", "te", "design",
		"of", "an", "automated", "mathematical", "assistant", "Some", "requirements", "for",
		"such", "a", "program", "are", "it", "must", "be", "easy", "to", "access", "and",
		"that", "the", "result", "must", "be", "obtained", "in", "a", "reasonably", "short",
		"time", "Accordingly", "the", "program", "is", "written", "for", "a", "time",
		"shared", "computer", "The", "Q", "32", "computer", "as", "System", "Development",
		"Corporation", "Santa", "Monica", "California", "was", "chosen", "because", "it",
		"also", "had", "a", "LISP", "1", "5", "compiler", "Programming", "and", "debugging",
		"was", "done", "from", "a", "remote", "teletype", "console", "at", "Stanford", "University",
	}
	filters := map[string]struct{}{"on": struct{}{}, "a": struct{}{}, "an": struct{}{}, "for": struct{}{},
		"of": struct{}{}, "it": struct{}{}, "to": struct{}{}}
	for i := 0; i < b.N; i++ {
		cwFilter(words, filters)
	}
}
