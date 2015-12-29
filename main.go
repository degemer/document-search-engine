package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Starting")
	r := CacmReader{path: "cacm/cacm.all"}

	for s := range r.Read() {
		fmt.Println(s)
	}
}
