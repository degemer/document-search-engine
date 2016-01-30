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

const VERSION = "0.1.0"

func main() {
	app := cli.NewApp()
	app.Name = "document-search-engine"
	app.Usage = "Search database"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "index, i",
			Value: "tf-idf",
			Usage: "Index to be used",
		},
		cli.StringFlag{
			Name:  "cacm",
			Value: "cacm",
			Usage: "Path to cacm directory",
		},
		cli.StringFlag{
			Name:  "save",
			Value: ".",
			Usage: "Path to indexes save directory",
		},
	}
	app.Commands = append([]cli.Command{}, index_command(), search_command(), measure_command())
	app.Version = VERSION
	app.Run(os.Args)
}

func index_command() cli.Command {
	return cli.Command{
		Name:    "index",
		Aliases: []string{"i"},
		Usage:   "Create index",
		Action: func(c *cli.Context) {
			options := make(map[string]string)
			options["cacm"] = c.GlobalString("cacm")
			options["saveDirectory"] = c.GlobalString("save")
			fmt.Println("Creating index", c.GlobalString("index"))
			i := index.New(c.GlobalString("index"), options)
			timeIndex(i.Create, "creation")
			timeIndex(i.Save, "save")
		},
	}
}

func search_command() cli.Command {
	c := cli.Command{
		Name:    "search",
		Aliases: []string{"s"},
		Usage:   "Search",
		Action: func(c *cli.Context) {
			options := make(map[string]string)
			options["cacm"] = c.GlobalString("cacm")
			options["saveDirectory"] = c.GlobalString("save")
			nbResults := c.Int("n")

			fmt.Println("Loading index", c.GlobalString("index"))
			i := index.New(c.GlobalString("index"), options)
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
				nbResultsPrinted := nbResults
				if len(results) < nbResultsPrinted {
					nbResultsPrinted = len(results)
				}
				fmt.Println("Search took", time.Since(start_search), "- printing", nbResultsPrinted, "out of", len(results))

				for i, res := range results {
					if i < nbResults {
						fmt.Println("Res", i, ":", res.Id, "-", res.Score)
					}
				}
				fmt.Println()
				fmt.Println("Enter request (empty request to exit): ")
				req, _ = reader.ReadString('\n')
			}
		},
	}
	c.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "n",
			Value: 10,
			Usage: "Number of results returned by request",
		},
	}
	return c
}

func measure_command() cli.Command {
	c := cli.Command{
		Name:    "measure",
		Aliases: []string{"m"},
		Usage:   "Measure search pertinence",
		Action: func(c *cli.Context) {
			options := make(map[string]string)
			options["cacm"] = c.GlobalString("cacm")
			options["saveDirectory"] = c.GlobalString("save")
			options["alpha"] = c.String("alpha")
			options["beta"] = c.String("beta")

			fmt.Println("Loading index", c.GlobalString("index"))
			i := index.New(c.GlobalString("index"), options)
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
	}
	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "alpha",
			Value: "0.5",
			Usage: "Alpha value for E-Measure",
		},
		cli.StringFlag{
			Name:  "beta",
			Value: "1",
			Usage: "Beta value for F-Measure",
		},
	}
	return c
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
