package bookstore

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/gin-gonic/gin"
)
var db *badger.DB


func Run() error {
	var err error
	db, err = badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		return err
	}
	defer db.Close()
	//clear the db using the clearDB function
	// err = clearDB()
	// if err != nil {
	// 	return err
	// }
	// //populate the db using the populateDB function
	// books := []Book{
	// 	{ID: "1", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 3},
	// 	{ID: "2", Title: "To Kill a Mocking",Quantity: 4, Author: "Harper Lee"},
	// 	{ID: "3", Title: "1984", Author: "George Orwell", Quantity: 5},
	// 	{ID: "4", Title: "Pride and Prejudice", Author: "Jane Austen", Quantity: 2},
	// }
	// err = populateDB(books)
	// if err != nil {
	// 	return err
	// }
	
	router := gin.Default()

	router.POST("/books", createBook)
	router.GET("/books/:id", getBookByID)
	router.GET("/books", getBooks)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	//insert a book
	router.POST("/insert", insertBook)
	err = router.Run("localhost:3000")
	if err != nil {
		return err
	}

	return nil
}