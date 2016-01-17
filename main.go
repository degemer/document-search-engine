package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/degemer/document-search-engine/index"
	"github.com/degemer/document-search-engine/measure"
	"github.com/degemer/document-search-engine/search"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	index_name := "tf-idf"
	cacm_path := "cacm"
	app.Name = "document-search-engine"
	app.Usage = "Search database"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "index, i",
			Value:       "tf-idf",
			Usage:       "Index to be used",
			Destination: &index_name,
		},
		cli.StringFlag{
			Name:        "cacm",
			Value:       "cacm",
			Usage:       "Path to cacm directory",
			Destination: &cacm_path,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "index",
			Aliases: []string{"i"},
			Usage:   "Create index",
			Action: func(c *cli.Context) {
				options := make(map[string]string)
				options["cacm"] = cacm_path
				fmt.Println("Creating index", index_name)
				i := index.New(index_name, options)
				timeIndex(i.Create, "creation")
				timeIndex(i.Save, "save")
			},
		},
		{
			Name:    "measure",
			Aliases: []string{"m"},
			Usage:   "Measure search pertinence",
			Action: func(c *cli.Context) {
				options := make(map[string]string)
				options["cacm"] = cacm_path

				fmt.Println("Loading index", index_name)
				i := index.New(index_name, options)
				loadOrCreateIndex(i)

				search_name := "vectorial"
				if len(c.Args()) > 0 {
					search_name = c.Args()[0]
				}
				fmt.Println("Creating", search_name, "search")
				searcher := search.New(search_name, i)

				measurer := measure.New("cacm", options)
				measurer.Measure(searcher)
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search",
			Action: func(c *cli.Context) {
				options := make(map[string]string)
				options["cacm"] = cacm_path

				fmt.Println("Loading index", index_name)
				i := index.New(index_name, options)
				loadOrCreateIndex(i)

				search_name := "vectorial"
				if len(c.Args()) > 0 {
					search_name = c.Args()[0]
				}
				fmt.Println("Creating", search_name, "search")
				searcher := search.New(search_name, i)

				reader := bufio.NewReader(os.Stdin)
				fmt.Println("Enter request: ")
				req, _ := reader.ReadString('\n')
				for req != "\n" {
					start_search := time.Now()
					results := searcher.Search(req)
					fmt.Println("Search took", time.Since(start_search), "printing 10 out of", len(results))

					for i, res := range results {
						if i < 10 {
							fmt.Println("Res", i, ":", res.Id, "-", res.Score)
						}
					}
					fmt.Println()
					fmt.Println("Enter request: ")
					req, _ = reader.ReadString('\n')
				}
			},
		},
	}
	app.Run(os.Args)
}

func timeIndex(fn func(), name string) {
	start := time.Now()
	fn()
	fmt.Println("Index", name, "took", time.Since(start))
}

func timeIndexError(fn func() error, name string) (err error) {
	start := time.Now()
	err = fn()
	if err != nil {
		return
	}
	fmt.Println("Index", name, "took", time.Since(start))
	return
}

func loadOrCreateIndex(i index.Index) {
	if err := timeIndexError(i.Load, "load"); err != nil {
		fmt.Println("Index not found, creating it")
		timeIndex(i.Create, "creation")
		fmt.Println("Saving it for future use")
		timeIndex(i.Save, "save")
	}
}
