package handler

import (
	"github.com/gin-gonic/gin"
	"library"
	"net/http"
	"strconv"
)

func (h *Handler) createAuthor(c *gin.Context) {
	var input library.Author
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.Author.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}

type getAllAuthorsResponse struct {
	Data []library.Author `json:"data"`
}

func (h *Handler) getAllAuthors(c *gin.Context) {

	authors, err := h.services.Author.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllAuthorsResponse{
		Data: authors,
	})

}

func (h *Handler) getAuthorById(c *gin.Context) {
	authorId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid author id param")
		return
	}

	author, err := h.services.Author.GetByID(authorId)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, author)
}

func (h *Handler) updateAuthor(c *gin.Context) {

	authorId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid author id param")
		return
	}

	var input library.UpdateAuthor
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.services.Author.Update(authorId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteAuthor(c *gin.Context) {
	authorId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid author id param")
		return
	}

	err = h.services.Author.Delete(authorId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
