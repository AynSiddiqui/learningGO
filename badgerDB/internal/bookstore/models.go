package bookstore

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger/v3"
)

type Book struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Author  string `json:"author"`
    Quantity int    `json:"quantity"`
}


func init() {
    var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}


//db commands:

//clearing all the data in the database
func clearDB() error {
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		//why is the prefetch size 10
		//PrefetchSize is the number of KV pairs to prefetch during iteration. This is useful to speed up iteration when the value sizes are small.
		//The default value is 100. If you know that your value sizes are small, you can reduce this value to reduce the number of disk reads.
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		// Iterate over all keys and delete them
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := txn.Delete(item.Key())
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func populateDB(books []Book) error {
	err := db.Update(func(txn *badger.Txn) error {
		for _, book := range books {
			// Check if any field is empty
			if book.ID == "" || book.Title == "" || book.Author == "" || book.Quantity <= 0 {
				fmt.Printf("Skipping book with ID %s: Incomplete data\n", book.ID)
				continue
			}

			// Marshal the book data to JSON
			bookJSON, err := json.Marshal(book)
			if err != nil {
				return err
			}

			// Set the book data in the database
			err = txn.Set([]byte(book.ID), bookJSON)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func displayAllBooks() {
	var books []Book
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			book := Book{}
			err := item.Value(func(val []byte) error {
				if err := json.Unmarshal(val, &book); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				log.Println("Error reading book:", err)
				continue
			}
			books = append(books, book)
		}
		return nil
	})

	if err != nil {
		log.Println("Error retrieving books:", err)
		return
	}

	log.Println("All Books:")
	for _, book := range books {
		log.Printf("ID: %s, Title: %s, Author: %s, Quantity: %d\n", book.ID, book.Title, book.Author, book.Quantity)
	}
}

