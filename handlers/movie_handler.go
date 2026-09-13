package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/models"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct { //The methods in your service are attached to the MovieService type, so you need a MovieService object to call them.and samre for other structs too
	service *services.MovieService
}

// constructor
func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

// create movie
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.Movie

	//read json request body(maps requested struct to backend struct)..client thing.they could send a bad request
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//call service...cuz its only handles not do core backend stuff... server thing..err != nil → Checks whether an error exists....err.Error() → Gets the text/message of that same error.
	if err := h.service.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{ //error for nt
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// get movie by its local primary key
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

// get movie by its TMDB ID
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

// Get all movies
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

// update movie
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid movie ID",
		})
		return
	}

	var movie models.Movie //create tempo object  this had everything unassigned or ID=0

	if err := c.ShouldBindJSON(&movie); err != nil { //this movie is updated one like the client request this new struct with changes..this i swhat client want it to be so we need to check it too and use bindjson
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	movie.ID = uint(id) //since new struct tempo one has ID=0 so assigning the one need to be

	if err := h.service.UpdateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, movie)

}

// delete movie
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
