package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStandardTokenise(t *testing.T) {
	content := `"An On-Line Program for Non-Numerical Algebra
	The goal of this program is to make a step toward te design of an automated mathematical assistant.
	Some requirements for such a program are: it must be easy to access,
	and that the result must be obtained in a reasonably short time.
	Accordingly the program is written for a time-shared computer.
	The Q-32 computer as System Development Corporation, Santa Monica, California,
	was chosen because it also had a LISP 1.5 compiler.
	Programming and debugging was done from a remote teletype console at Stanford University.`
	result_words := []string{
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
	rawDocument := RawDocument{Id: 35, Content: strings.Replace(content, "\n", "", -1)}
	tokenizedDocument := StandardTokenize(rawDocument)

	assert.Equal(t, tokenizedDocument.Words, result_words, "They should be equal")
	assert.Equal(t, tokenizedDocument.Id, 35, "They should be equal")
}

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
	tokenizedDocument := TokenizedDocument{Id: 45, Words: words}
	filteredDocument := CWFilter(tokenizedDocument, filters)

	assert.Equal(t, filteredDocument.Words, filtered_words, "Filter incorrect")
	assert.Equal(t, filteredDocument.Id, 45, "Id modified...")
}

func TestCountWord(t *testing.T) {
	words := []string{"the", "an", "line", "the", "program", "for", "for", "non", "numerical", "the"}
	countedWords := map[string]int{"the": 3, "an": 1, "line": 1, "program": 1, "for": 2, "non": 1, "numerical": 1}
	assert.Equal(t, CountWords(words), countedWords, "Count incorrect")
}
