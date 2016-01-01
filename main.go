package main

import (
	// "fmt"
	"github.com/degemer/document-search-engine/index"
	"log"
	"time"
)

func main() {
	log.Println("Starting")
	start := time.Now()
	i := index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"})
	i.Create()
	elapsed := time.Since(start)
	log.Printf("Index creation took %s", elapsed)
}
