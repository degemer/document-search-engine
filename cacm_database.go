package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CacmReader struct {
	path string
}

func (c CacmReader) Read() <-chan RawDocument {
	return c.parseDocument(c.scanDatabase())
}

func (c CacmReader) scanDatabase() <-chan string {
	file, err := os.Open(c.path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(ScanCacmDocument)
	unparsed_strings := make(chan string)
	go func() {
		for scanner.Scan() {
			unparsed_strings <- scanner.Text()
		}
		// Can't be deferred, otherwise is closed before the scan starts
		file.Close()

		if err := scanner.Err(); err != nil {
			log.Fatalln("Error scanning file: ", err)
		}
		close(unparsed_strings)
	}()
	return unparsed_strings
}

func (c CacmReader) parseDocument(unparsed_strings <-chan string) <-chan RawDocument {
	regex_expr := `(?ms)^\.I (\d+).^\.T.(.+?)(?:.^\.W.(.+?))?.^\.B.(.+?)(?:.^\.A.(.+?))?(?:.^\.K.(.+?))?(?:.^\.C.(.+?))?.^\.N.(.+?).^\.X.(.+?)`
	r := regexp.MustCompile(regex_expr)
	raw_documents := make(chan RawDocument)
	go func() {
		for s := range unparsed_strings {
			regex_result := r.FindAllStringSubmatch(s, -1)[0]
			I, err := strconv.Atoi(regex_result[1])
			if err != nil {
				log.Fatalln("Unable to convert id ", regex_result[1], "to int: ", err)
			}
			raw_documents <- RawDocument{Id: I, Content: strings.Replace(regex_result[2]+" "+regex_result[3]+" "+regex_result[6], "\n", " ", -1)}
		}
		close(raw_documents)
	}()
	return raw_documents
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
