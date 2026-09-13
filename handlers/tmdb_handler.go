package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type TMDBHandler struct {
	service               *services.TMDBService
	recommendationService *services.RecommendationService
}

// constructor
func NewTMDBHandler(service *services.TMDBService, recommendationService *services.RecommendationService) *TMDBHandler {
	return &TMDBHandler{
		service:               service,
		recommendationService: recommendationService,
	}
}

// searchMovies searches TMDb by movie title.
// example:GET /tmdb/search?query=interstellar
func (h *TMDBHandler) SearchMovies(c *gin.Context) {

	//get movie title from query parameter
	query := c.Query("query")

	//validate that a search a search query is validated
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter is required",
		})
		return
	}

	//call the service to search movies on TMDB
	result, err := h.service.SearchMovies(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the search results as JSON
	c.JSON(http.StatusOK, result)
}

// getMovieDetails returns details of a single TMDb movie.
// example:GET /tmdb/movie/157336
func (h *TMDBHandler) GetMovieDetails(c *gin.Context) {

	// Convert the movie ID from the URL into an integer
	id, err := strconv.Atoi(c.Param("id")) //atoi converts string to int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid movie id",
		})
		return
	}

	// Fetch movie details from the service
	movie, err := h.service.GetMovieDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return the movie details as JSON
	c.JSON(http.StatusOK, movie)
}

// publish a tmdb movie to local database
func (h *TMDBHandler) PublishMovie(c *gin.Context) {

	// get tmdb movie id from URL
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid TMDb movie ID",
		})
		return
	}

	//publish the movie
	movie, err := h.service.PublishMovie(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.recommendationService.GenerateAndStoreMovieEmbedding(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	// return the published movie
	c.JSON(http.StatusCreated, movie)
}
