package main

import (
	"github.com/codegangsta/cli"
	"github.com/degemer/document-search-engine/index"
	"github.com/degemer/document-search-engine/measure"
	"github.com/degemer/document-search-engine/search"
	"log"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "document-search-engine"
	app.Usage = "Search database"
	app.Action = func(c *cli.Context) {
		log.Println("Starting")
		start_creation := time.Now()
		i := index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"})
		i.Create()
		elapsed_creation := time.Since(start_creation)
		log.Printf("Index creation took %s", elapsed_creation)
		start_save := time.Now()
		i.Save()
		elapsed_save := time.Since(start_save)
		log.Printf("Index save took %s", elapsed_save)

		start_load := time.Now()
		j := index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"})
		j.Load()
		elapsed_load := time.Since(start_load)
		log.Printf("Index load took %s", elapsed_load)

		searcher := search.New("vectorial", i)
		start_measure := time.Now()
		measurer := measure.New("cacm", map[string]string{"cacm_path": "cacm"})
		measurer.Measure(searcher)
		elapsed_measure := time.Since(start_measure)
		log.Printf("Measure took %s", elapsed_measure)
	}
	app.Run(os.Args)
}
