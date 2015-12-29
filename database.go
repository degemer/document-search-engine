package main

type Reader interface {
	Read() <-chan RawDocument
}

type RawDocument struct {
	Id      int
	Content string
}
