package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type UserGenreHandler struct {
	service *services.UserGenreService
}

// constructor
func NewUserGenreHandler(service *services.UserGenreService) *UserGenreHandler {
	return &UserGenreHandler{
		service: service,
	}
}

// GetFavoriteGenres returns the authenticated user's favorite genres.
func (h *UserGenreHandler) GetFavoriteGenres(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	genres, err := h.service.GetFavoriteGenres(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, genres)
}

// AddFavoriteGenre adds a genre to the authenticated user's favorite genres.
func (h *UserGenreHandler) AddFavoriteGenre(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	tmdbGenreID, err := strconv.ParseUint(c.Param("genreID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid genre ID",
		})
		return
	}

	err = h.service.AddFavoriteGenre(userID, int(tmdbGenreID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Genre added to favorite genres successfully",
	})
}

// RemoveFavoriteGenre removes a genre from the authenticated user's favorites.
func (h *UserGenreHandler) RemoveFavoriteGenre(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	tmdbGenreID, err := strconv.ParseUint(c.Param("genreID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid genre ID",
		})
		return
	}

	err = h.service.RemoveFavoriteGenre(userID, int(tmdbGenreID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Genre removed from favorite genres successfully",
	})
}
