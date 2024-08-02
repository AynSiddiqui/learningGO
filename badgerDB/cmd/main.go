package main

import (
	"badgerdb/internal/bookstore"
	"log"
)

func main() {
	err := bookstore.Run()
	if err != nil {
		log.Fatal(err)
	}
}
