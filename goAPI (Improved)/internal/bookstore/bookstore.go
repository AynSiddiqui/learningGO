package bookstore

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getBooks(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
    var newBook Book
    if err := c.BindJSON(&newBook); err != nil {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    books = append(books, newBook)
    c.IndentedJSON(http.StatusCreated, newBook)
}

func bookById(c *gin.Context) {
    id := c.Param("id")
    book, err := getBookById(id)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
        return
    }
    c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
    id, ok := c.GetQuery("id")
    if !ok {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
        return
    }
    book, err := getBookById(id)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
        return
    }
    if book.Quantity <= 0 {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not available"})
        return
    }
    book.Quantity -= 1
    c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
    id, ok := c.GetQuery("id")
    if !ok {
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "id is required"})
        return
    }
    book, err := getBookById(id)
    if err != nil {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found"})
        return
    }
    book.Quantity += 1
    c.IndentedJSON(http.StatusOK, book)
}

func RunServer() {
    // Ensure books are loaded from file before starting the server
    if err := LoadBooksFromFile(); err != nil {
        panic(err) // Handle error appropriately
    }
	fmt.Println("Loaded books from file", books)
    router := gin.Default()
    router.GET("/books", getBooks)
    router.POST("/books", createBook)
    router.GET("/books/:id", bookById)
    router.PATCH("/checkout", checkoutBook)
    router.PATCH("/return", returnBook)
    router.Run("localhost:3000")
}
