package main

import (
	"fmt"
	"github.com/degemer/document-search-engine/index"
	"log"
	"time"
)

func main() {
	log.Println("Starting")
	reader := index.CacmReader{Path: "cacm/cacm.all"}
	tokenizer := index.StandardTokenizer{}
	filter := index.CommonWordsFilter{Path: "cacm/common_words"}
	counter := index.StandardCounter{}

	start := time.Now()
	for s := range index.Tf(counter.Count(filter.Filter(tokenizer.Tokenize(reader.Read())))) {
		if s.Id == 3204 {
			fmt.Println(s)
		}
		// fmt.Println(s)
	}
	elapsed := time.Since(start)
	log.Printf("Index creation took %s", elapsed)
}
