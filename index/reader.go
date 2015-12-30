package index

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
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
	regex_expr := `(?ms)^\.I (\d+).^\.T.(.+?)(?:.^\.W.(.+?))?.^\.B.(.+?)(?:.^\.A.(.+?))?(?:.^\.K.(.+?))?(?:.^\.C.(.+?))?.^\.N.(.+?).^\.X.(.+?)`
	r := regexp.MustCompile(regex_expr)
	rawDocuments := make(chan RawDocument)
	go func() {
		for s := range unparsedStrings {
			regex_result := r.FindAllStringSubmatch(s, -1)[0]
			I, err := strconv.Atoi(regex_result[1])
			if err != nil {
				log.Fatalln("Unable to convert id ", regex_result[1], "to int: ", err)
			}
			rawDocuments <- RawDocument{Id: I, Content: strings.Replace(regex_result[2]+" "+regex_result[3]+" "+regex_result[6], "\n", " ", -1)}
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
