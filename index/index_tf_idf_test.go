package index

import (
	"testing"
)

func BenchmarkTfIdfCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-idf", options)
		i.Create()
	}
}

func BenchmarkTfIdfStemCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-idf-stem", options)
		i.Create()
	}
}

func BenchmarkTfIdfNormCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-idf-norm", options)
		i.Create()
	}
}

func BenchmarkTfIdfNormStemCreate(b *testing.B) {
	options := make(map[string]string)
	options["cacm"] = "../cacm"
	for ind := 0; ind < b.N; ind++ {
		i := New("tf-idf-norm-stem", options)
		i.Create()
	}
}
