package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"regexp"
	"strings"
)

type CacmReader struct {
	path             string
	unparsed_strings chan RawDocument
}

type CacmDocument struct {
	I string
	T string
	W string
	K string
}

func (c CacmReader) Read() <-chan RawDocument {
	file, err := os.Open(c.path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(ScanCacmDocument)
	c.unparsed_strings = make(chan RawDocument)

	go func() {
		regex_I := regexp.MustCompile(`.I (\d+)`)
		regex_T := regexp.MustCompile(` \.T (.+?) \.[A-Z]`)
		regex_W := regexp.MustCompile(` \.W (.+?) \.[A-Z]`)
		regex_K := regexp.MustCompile(` \.K (.+?) \.[A-Z]`)
		for scanner.Scan() {
			doc := &CacmDocument{}
			s := strings.Replace(scanner.Text(), "\n", " ", -1)
			doc.I = findFirstSubmatch(regex_I, s)
			if doc.I == "" {
				log.Println(s)
			}
			doc.T = findFirstSubmatch(regex_T, s)
			doc.W = findFirstSubmatch(regex_W, s)
			doc.K = findFirstSubmatch(regex_K, s)
			c.unparsed_strings <- RawDocument{Id: doc.I, Content: doc.T + " " + doc.W + " " + doc.K}
		}
		// Can't be deferred, otherwise is closed before the scan starts
		file.Close()

		if err := scanner.Err(); err != nil {
			log.Fatalln("Error scanning file: ", err)
		}
		close(c.unparsed_strings)
	}()
	return c.unparsed_strings
}

func ScanCacmDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data[1:], []byte(".I")); i >= 0 {
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

func findFirstSubmatch(r *regexp.Regexp, s string) string {
	s_slice := r.FindStringSubmatch(s)
	if s_slice == nil {
		return ""
	}
	return s_slice[1]
}
