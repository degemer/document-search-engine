package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.Println("Starting")
	reader := CacmReader{path: "cacm/cacm.all"}
	tokenizer := StandardTokenizer{}
	filter := CommonWordsFilter{path: "cacm/common_words"}
	counter := StandardCounter{}

	start := time.Now()
	for s := range counter.Count(filter.Filter(tokenizer.Tokenize(reader.Read()))) {
		if s.Id == 1 {
			fmt.Println(s)
		}
		// fmt.Println(s)
	}
	elapsed := time.Since(start)
	log.Printf("Index creation took %s", elapsed)
}
