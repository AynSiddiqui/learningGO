package bookstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)

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



//simple create book below but dangerous because it does not check if the book already exists and overwrites the book
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


//insert the book, and dynamically allocate the id after the last id in case the book name and all fields are not similar with any previous book
// first check if the book already exists and the name is the same and the author is the same
// if it is, update the quantity
// if it is not, get the last id and increment it by 1
// then insert the book
// if the book already exists but the name is the same and the author is different, return an error
// if the book does not exist, insert the book
// write the web server code to insert a book using a helper function, input will be a json object with the book name and author and quantity, check if the book already exists, if it does, update the quantity, if it does not, insert the book with the next id
func insertBook(c *gin.Context) {
	var newBook Book
	if err := c.BindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var book Book
	var id string
	err := db.Update(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
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
			if book.Title == newBook.Title && book.Author == newBook.Author {
				book.Quantity += newBook.Quantity
				bookJSON, err := json.Marshal(book)
				if err != nil {
					return err
				}
				return txn.Set([]byte(book.ID), bookJSON)
			}
		}
		//get the last id and do the increment
		it.Rewind()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			id = string(item.Key())
		}
		//increment the id
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
		idInt++
		newBook.ID = strconv.Itoa(idInt)
		bookJSON, err := json.Marshal(newBook)
		if err != nil {
			return err
		}
		return txn.Set([]byte(newBook.ID), bookJSON)
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