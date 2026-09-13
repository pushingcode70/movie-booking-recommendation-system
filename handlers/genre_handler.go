package handlers

import (
	"net/http"

	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	service *services.GenreService
}

// constructor
func NewGenreHandler(service *services.GenreService) *GenreHandler {
	return &GenreHandler{
		service: service,
	}
}

// get all genres
func (h *GenreHandler) GetAllGenres(c *gin.Context) {
	genres, err := h.service.GetAllGenres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, genres)
}
