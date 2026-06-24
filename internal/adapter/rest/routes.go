package rest

import (
	"go-api-boilerplate/internal/adapter/rest/book"
	"go-api-boilerplate/internal/adapter/rest/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine, bookHandler *book.BookHandler) {
	router.Use(middlewares.ErrorHandler())

	booksGroup := router.Group("/books")
	{
		booksGroup.POST("", bookHandler.CreateBook)
		booksGroup.GET("/:id", bookHandler.GetBook)
		booksGroup.GET("", bookHandler.GetBooks)
		booksGroup.PUT("/:id", bookHandler.UpdateBook)
		booksGroup.DELETE("/:id", bookHandler.DeleteBook)
	}
}
