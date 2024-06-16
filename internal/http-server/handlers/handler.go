package handler

import (
	"github.com/gin-gonic/gin"
	"library/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	api := router.Group("/api")
	{

		authors := api.Group("/authors")
		{

			authors.POST("/", h.createAuthor)
			authors.GET("/", h.getAllAuthors)
			authors.GET("/:id", h.getAuthorById)
			authors.PUT("/:id", h.updateAuthor)
			authors.DELETE("/:id", h.deleteAuthor)

		}

		books := api.Group("/books")
		{

			books.POST("/", h.createBook)
			books.GET("/", h.getAllBooks)
			books.GET("/:id", h.getBookById)
			books.PUT("/:id", h.updateBook)
			books.DELETE("/:id", h.deleteBook)

		}

		api.PUT("book/:book_id/authors/:author_id", h.updateBookAndAuthor)
	}

	return router

}
