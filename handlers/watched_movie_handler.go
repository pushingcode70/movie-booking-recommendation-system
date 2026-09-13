package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/dto"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type WatchedMovieHandler struct {
	service *services.WatchedMovieService
}

func NewWatchedMovieHandler(service *services.WatchedMovieService) *WatchedMovieHandler {
	return &WatchedMovieHandler{
		service: service,
	}
}

func (h *WatchedMovieHandler) AddWatchedMovie(c *gin.Context) {
	var req dto.CreateWatchedMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if err := h.service.AddWatchedMovie(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Movie marked as watched"})
}

func (h *WatchedMovieHandler) UpdateWatchedMovie(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdbID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TMDb ID"})
		return
	}

	var req dto.UpdateWatchedMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if err := h.service.UpdateWatchedMovie(userID, tmdbID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Watched movie review updated"})
}

func (h *WatchedMovieHandler) GetWatchedMovies(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	response, err := h.service.GetWatchedMovies(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *WatchedMovieHandler) RemoveWatchedMovie(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	tmdbID, err := strconv.Atoi(c.Param("tmdbID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TMDb ID"})
		return
	}

	if err := h.service.RemoveWatchedMovie(userID, tmdbID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Movie removed from watched list"})
}
