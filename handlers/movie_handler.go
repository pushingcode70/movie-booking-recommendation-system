package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/models"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	service *services.MovieService
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.Movie

	// read json request body
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// call service to create movie
	if err := h.service.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	movie, err := h.service.GetMovieByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetMovieByTMDBID(c *gin.Context) {
	tmdbID, err := strconv.Atoi(c.Param("tmdbId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid TMDB ID"})
		return
	}

	movie, err := h.service.GetMovieByTMDBID(tmdbID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.service.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	var movie models.Movie

	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	movie.ID = uint(id)

	if err := h.service.UpdateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)

}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	if err := h.service.DeleteMovie(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully",
	})
}
