package index

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

	assert.Equal(t, standardTokenize(content), result_words, "They should be equal")
}

func BenchmarkStandardTokenize(b *testing.B) {
	content := strings.Replace(`"An On-Line Program for Non-Numerical Algebra
 The goal of this program is to make a step toward te design of an automated mathematical assistant.
 Some requirements for such a program are: it must be easy to access,
 and that the result must be obtained in a reasonably short time.
 Accordingly the program is written for a time-shared computer.
 The Q-32 computer as System Development Corporation, Santa Monica, California,
 was chosen because it also had a LISP 1.5 compiler.
 Programming and debugging was done from a remote teletype console at Stanford University.`,
		"\n", "", -1)
	for i := 0; i < b.N; i++ {
		standardTokenize(content)
	}
}
