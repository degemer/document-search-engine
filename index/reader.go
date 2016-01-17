package index

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Reader interface {
	Read() <-chan RawDocument
	ReadOne(string) RawDocument
}

type RawDocument struct {
	Id      int
	Content string
}

type CacmReader struct {
	Path string
}

func NewReader(options map[string]string) Reader {
	if options["cacm_file"] != "" {
		return CacmReader{Path: filepath.Join(options["cacm_file"])}
	}
	return CacmReader{Path: filepath.Join(options["cacm"], "cacm.all")}
}

func (c CacmReader) Read() <-chan RawDocument {
	return c.parseDocument(c.scanDatabase())
}

func (c CacmReader) scanDatabase() <-chan string {
	file, err := os.Open(c.Path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(scanCacmDocument)
	unparsedStrings := make(chan string, CHANNEL_SIZE)
	go func() {
		for scanner.Scan() {
			unparsedStrings <- scanner.Text()
		}
		// Can't be deferred, otherwise is closed before the scan starts
		file.Close()

		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file: ", err)
			os.Exit(1)
		}
		close(unparsedStrings)
	}()
	return unparsedStrings
}

func (c CacmReader) parseDocument(unparsedStrings <-chan string) <-chan RawDocument {
	rawDocuments := make(chan RawDocument, CHANNEL_SIZE)
	go func() {
		for s := range unparsedStrings {
			rawDocuments <- c.ReadOne(s)
		}
		close(rawDocuments)
	}()
	return rawDocuments
}

func (c CacmReader) ReadOne(doc string) RawDocument {
	id, content := cacmDoc(doc)
	return RawDocument{Id: id, Content: content}
}

func scanCacmDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte("\n.I")); i >= 0 {
		// We have a full doc
		return i + 1, data[0:i], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func cacmDoc(doc string) (int, string) {
	baseValues := []string{"T", "W", "B", "A", "K", "C", "N", "X"}
	presentValues := []string{}
	indValues := []int{}
	content := ""

	for _, val := range baseValues {
		if ind := strings.Index(doc, "\n."+val); ind != -1 {
			presentValues = append(presentValues, val)
			indValues = append(indValues, ind)
		}
	}
	id, err := strconv.Atoi(doc[3:indValues[0]])
	if err != nil {
		fmt.Println("Unable to convert id ", doc[3:indValues[0]], "to int: ", err)
		os.Exit(1)
	}

	for i, val := range presentValues {
		if val == "T" || val == "W" || val == "K" {
			content += doc[indValues[i]+3 : indValues[i+1]]
		}
	}
	return id, content
}
