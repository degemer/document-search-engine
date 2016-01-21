package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPorter(t *testing.T) {
	words := []string{
		"an", "line", "program", "for", "non", "numerical", "the", "goal", "of",
		"this", "program", "is", "to", "make", "step", "toward", "te",
		"design", "of", "an", "automated", "mathematical", "assistant",
		"some", "requirements", "for",
	}
	stemmedWords := []string{
		"an", "line", "program", "for", "non", "numer", "the", "goal", "of",
		"thi", "program", "is", "to", "make", "step", "toward", "te",
		"design", "of", "an", "autom", "mathemat", "assist",
		"some", "requir", "for",
	}

	assert.Equal(t, porter(words), stemmedWords, "Stemmer incorrect")
}

func BenchmarkPorter(b *testing.B) {
	words := []string{
		"an", "line", "program", "for", "non", "numerical", "the", "goal", "of",
		"this", "program", "is", "to", "make", "step", "toward", "te",
		"design", "of", "an", "automated", "mathematical", "assistant",
		"some", "requirements", "for",
	}
	for i := 0; i < b.N; i++ {
		porter(words)
	}
}
