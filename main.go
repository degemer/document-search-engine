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

	start := time.Now()
	for s := range tokenizer.Tokenize(reader.Read()) {
		fmt.Println(s)
	}
	elapsed := time.Since(start)
	log.Printf("Index creation took %s", elapsed)
}
