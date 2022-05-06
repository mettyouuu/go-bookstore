package main

import (
	"errors"
	_ "errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []Book{
	{ID: "1", Title: "Holy", Author: "Kopi", Quantity: 12},
	{ID: "2", Title: "Devil", Author: "Kopi", Quantity: 2},
	{ID: "3", Title: "English", Author: "Kopi", Quantity: 8},
	{ID: "4", Title: "Happy", Author: "Kopi", Quantity: 9},
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Book with id " + id + " not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadGateway, gin.H{"Message": "Book not available"})
		return
	} else {
		book.Quantity -= 1
		c.IndentedJSON(http.StatusOK, book)
	}

}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Book with id " + id + " not found"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Book with id " + id + " not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*Book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}
	return nil, errors.New("book with id " + id + " not found")
}

func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books := append(books, newBook)
	c.IndentedJSON(http.StatusCreated, books)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)

	err := router.Run(":8080")
	if err != nil {
		return
	}

}
