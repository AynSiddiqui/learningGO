package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

type Book struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Quantity int    `json:"quantity"`
}

var db *badger.DB
var initialBooks = []Book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func main() {
	var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Clear the database
	// err = clearDB()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Populate the database with initial books
	// err = populateDB(initialBooks)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// displayAllBooks()


	router := gin.Default()

	// Define the routes
	router.POST("/books", createBook)
	router.GET("/books/:id", getBookByID)
	router.GET("/books", getBooks)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	
	// Run the server
	router.Run("localhost:3000")	
}


func clearDB() error {
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
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

			// Marshal the book data
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


func getBooks(c *gin.Context) {
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
				// Print the data before unmarshaling
				fmt.Println("Value from database:", string(val))

				// Check if the data is valid JSON
				if !json.Valid(val) {
					fmt.Println("Invalid JSON data:", string(val))
					return nil // Skip this data and continue to the next iteration
				}

				// Unmarshal the value into a Book struct
				if err := json.Unmarshal(val, &book); err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}
			fmt.Println("Parsed book:", book)
			books = append(books, book)
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}




func createBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := db.Update(func(txn *badger.Txn) error {
		bookJSON, err := json.Marshal(newBook)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(newBook.ID), bookJSON)
		return err
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newBook)
}

func getBookByID(c *gin.Context) {
	id := c.Param("id")

	var searchedBook Book
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		// Check if the item exists
		if item == nil {
			return errors.New("Book not found")
		}

		err = item.Value(func(val []byte) error {
			// Check if the data is valid JSON
			if !json.Valid(val) {
				return errors.New("invalid JSON data")
			}

			// Unmarshal the data into the searchedBook struct
			if err := json.Unmarshal(val, &searchedBook); err != nil {
				return err
			}

			// Check if the ID matches
			if searchedBook.ID != id {
				return errors.New("Book not found")
			}

			return nil
		})
		return err
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, searchedBook)
}


func checkoutBook(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var book Book
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &book)
		})
		if err != nil {
			return err
		}
		if book.Quantity <= 0 {
			return errors.New("book not available")
		}
		book.Quantity--
		bookJSON, err := json.Marshal(book)
		if err != nil {
			return err
		}
		return txn.Set([]byte(id), bookJSON)
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var book Book
	err := db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &book)
		})
		if err != nil {
			return err
		}
		book.Quantity++
		bookJSON, err := json.Marshal(book)
		if err != nil {
			return err
		}
		return txn.Set([]byte(id), bookJSON)
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}
