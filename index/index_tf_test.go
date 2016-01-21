package index

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWordsTfFrequency(t *testing.T) {
	wordsCount := map[string]int{
		"test": 1, "temp": 10, "aa": 100,
	}
	expectedTfFrequency := map[string]float64{
		"test": 1, "temp": 2, "aa": 3,
	}
	assert.Equal(t, wordsTfFrequency(wordsCount), expectedTfFrequency, "Incorrect frequency")
}

func BenchmarkTfNormCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-norm", options)
		i.Create()
	}
}

func BenchmarkTfNormStemCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-norm-stem", options)
		i.Create()
	}
}
