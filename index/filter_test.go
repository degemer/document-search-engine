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
