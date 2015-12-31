package index

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

type Reader interface {
	Read() <-chan RawDocument
}

type RawDocument struct {
	Id      int
	Content string
}

type CacmReader struct {
	Path string
}

func (c CacmReader) Read() <-chan RawDocument {
	return c.parseDocument(c.scanDatabase())
}

func (c CacmReader) scanDatabase() <-chan string {
	file, err := os.Open(c.Path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(ScanCacmDocument)
	unparsedStrings := make(chan string)
	go func() {
		for scanner.Scan() {
			unparsedStrings <- scanner.Text()
		}
		// Can't be deferred, otherwise is closed before the scan starts
		file.Close()

		if err := scanner.Err(); err != nil {
			log.Fatalln("Error scanning file: ", err)
		}
		close(unparsedStrings)
	}()
	return unparsedStrings
}

func (c CacmReader) parseDocument(unparsedStrings <-chan string) <-chan RawDocument {
	rawDocuments := make(chan RawDocument)
	go func() {
		for s := range unparsedStrings {
			id, content := parseDoc(s)
			rawDocuments <- RawDocument{Id: id, Content: content}
		}
		close(rawDocuments)
	}()
	return rawDocuments
}

func ScanCacmDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
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

func parseDoc(doc string) (int, string) {
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
		log.Fatalln("Unable to convert id ", doc[3:indValues[0]], "to int: ", err)
	}

	for i, val := range presentValues {
		if val == "T" || val == "W" || val == "K" {
			content += doc[indValues[i]+3 : indValues[i+1]]
		}
	}
	return id, content
}
