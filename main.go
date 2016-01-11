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

		temp := search.New("vectorial", i)
		req := `Preliminary Report-International Algebraic Language`
 		for i, res := range(temp.Search(req)) {
			log.Println("Res", i, ":", res.Id, "-", res.Score)
 		}
	}
	app.Run(os.Args)
}
