package main

import (
	"github.com/codegangsta/cli"
	"github.com/degemer/document-search-engine/index"
	"github.com/degemer/document-search-engine/search"
	"log"
	"os"
	"time"
)

func main() {
	temp := search.New("test", index.New("tf-idf", map[string]string{"cacm_path": "cacm/cacm.all", "common_words_path": "cacm/common_words"}))
	temp.Search("test")
	app := cli.NewApp()
	app.Name = "document-search-engine"
	app.Usage = "Search database"
	app.Action = func(c *cli.Context) {
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
	app.Run(os.Args)
}
