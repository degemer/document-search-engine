package parser

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestStandardTokenise(t *testing.T) {
	content := strings.Replace(`"An On-Line Program for Non-Numerical Algebra
 The goal of this program is to make a step toward te design of an automated mathematical assistant.
 Some requirements for such a program are: it must be easy to access,
 and that the result must be obtained in a reasonably short time.
 Accordingly the program is written for a time-shared computer.
 The Q-32 computer as System Development Corporation, Santa Monica, California,
 was chosen because it also had a LISP 1.5 compiler.
 Programming and debugging was done from a remote teletype console at Stanford University.`,
		"\n", "", -1)
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

	assert.Equal(t, StandardTokenize(content), result_words, "They should be equal")
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

	assert.Equal(t, CWFilter(words, filters), filtered_words, "Filter incorrect")
}

func TestCountWord(t *testing.T) {
	words := []string{"the", "an", "line", "the", "program", "for", "for", "non", "numerical", "the"}
	countedWords := map[string]int{"the": 3, "an": 1, "line": 1, "program": 1, "for": 2, "non": 1, "numerical": 1}
	assert.Equal(t, CountWords(words), countedWords, "Count incorrect")
}

func TestCacmDoc(t *testing.T) {
	unparsed_req := `.I 1
.W
 What articles exist which deal with TSS (Time Sharing System), an
operating system for IBM computers?
.N
 1. Richard Alexander, Comp Serv, Langmuir Lab (TSS)`
 	expected_req := `
 What articles exist which deal with TSS (Time Sharing System), an
operating system for IBM computers?`
	parsed_req_id, parsed_req := CacmDoc(unparsed_req)
	assert.Equal(t, parsed_req_id, 1, "Id incorrect")
	assert.Equal(t, parsed_req, expected_req, "Content incorrect")

	unparsed_doc := `.I 1655
.T
Code Extension Procedures for Information
Interchange* (Proposed USA Standard)
.B
CACM December, 1968
.K
standard code, code, information interchange, characters,
shift out, shift in, escape, data link
escape, control functions, standard procedures,
code extension, code table, bit pattern
.C
1.0 2.0 2.43 3.20 3.24 3.50 3.51 3.52 3.53 3.54 3.55 3.56 3.57 3.70 3.71 3.72
3.73, 3.74, 3.75, 3.80, 3.81, 3.82, 3.83, 5.0, 5.1, 6.2, 6.21, 6.22
.N
CA681211 JB February 21, 1978  12:16 PM
.X
1655	5	1655
1655	5	1655
1655	5	1655`
	expected_doc := `
Code Extension Procedures for Information
Interchange* (Proposed USA Standard)
standard code, code, information interchange, characters,
shift out, shift in, escape, data link
escape, control functions, standard procedures,
code extension, code table, bit pattern`
	parsed_doc_id, parsed_doc := CacmDoc(unparsed_doc)
	assert.Equal(t, parsed_doc_id, 1655, "Id incorrect")
	assert.Equal(t, parsed_doc, expected_doc, "Content incorrect")
}
