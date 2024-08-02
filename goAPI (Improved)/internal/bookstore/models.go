package bookstore

import (
	"encoding/json"
	"fmt"
	"os"
)

type Book struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Author  string `json:"author"`
    Quantity int    `json:"quantity"`
}

const filePath = "db/body.json"

var books []Book

func init() {
    if err := LoadBooksFromFile(); err != nil {
        // Handle error if necessary
		return
    }
}

func getBookById(id string) (*Book, error) {
    for i, b := range books {
        if b.ID == id {
            return &books[i], nil
        }
    }
    return nil, fmt.Errorf("book with ID %s not found", id)
}
func SaveBooksToFile() error {
    // Load existing books from the file
    err := LoadBooksFromFile()
    if err != nil {
        return err
    }

    // Write the updated books back to the file
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    if err := encoder.Encode(books); err != nil {
        return err
    }

    return nil
}

func LoadBooksFromFile() error {
    file, err := os.Open(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil
        }
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&books); err != nil {
        return err
    }

    // Print loaded books
    fmt.Printf("Loaded %d books from file\n", len(books))
    for _, b := range books {
        fmt.Println(b)
    }

    return nil
}