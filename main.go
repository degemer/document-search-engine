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

	start := time.Now()
	for s := range filter.Filter(tokenizer.Tokenize(reader.Read())) {
		fmt.Println(s)
	}
	elapsed := time.Since(start)
	log.Printf("Index creation took %s", elapsed)
}
