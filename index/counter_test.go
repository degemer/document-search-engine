package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCountWord(t *testing.T) {
	words := []string{"the", "an", "line", "the", "program", "for", "for", "non", "numerical", "the"}
	countedWords := map[string]int{"the": 3, "an": 1, "line": 1, "program": 1, "for": 2, "non": 1, "numerical": 1}
	assert.Equal(t, countWords(words), countedWords, "Count incorrect")
}

func BenchmarkCountWord(b *testing.B) {
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
	for i := 0; i < b.N; i++ {
		countWords(words)
	}
}
