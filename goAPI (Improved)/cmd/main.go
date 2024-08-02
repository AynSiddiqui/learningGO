package main

import (
	"goAPI/internal/bookstore"
)

func main() {
	bookstore.LoadBooksFromFile()
    bookstore.RunServer()
}
