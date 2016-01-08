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
