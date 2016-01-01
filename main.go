package main

import (
	"github.com/degemer/document-search-engine/index"
	"log"
	"time"
)

func main() {
	log.Println("Starting")
	start_creation := time.Now()
	i := index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"})
	i.Create()
	i.Save()
	elapsed_creation := time.Since(start_creation)
	log.Printf("Index creation took %s", elapsed_creation)

	start_load := time.Now()
	j := index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"})
	j.Load()
	elapsed_load := time.Since(start_load)
	log.Printf("Index load took %s", elapsed_load)
}
