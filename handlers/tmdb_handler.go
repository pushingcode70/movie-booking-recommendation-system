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

func NewTMDBHandler(service *services.TMDBService, recommendationService *services.RecommendationService) *TMDBHandler {
	return &TMDBHandler{
		service:               service,
		recommendationService: recommendationService,
	}
}

func (h *TMDBHandler) SearchMovies(c *gin.Context) {

	// get movie title from query parameter
	query := c.Query("query")

	// validate search query parameter
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter is required",
		})
		return
	}

	// call service to search movies on tmdb
	result, err := h.service.SearchMovies(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return search results as json
	c.JSON(http.StatusOK, result)
}

func (h *TMDBHandler) GetMovieDetails(c *gin.Context) {

	// convert movie id from url to integer
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid movie id",
		})
		return
	}

	// fetch movie details from service
	movie, err := h.service.GetMovieDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// return movie details as json
	c.JSON(http.StatusOK, movie)
}

func (h *TMDBHandler) PublishMovie(c *gin.Context) {

	// get tmdb movie id from url
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid TMDb movie ID",
		})
		return
	}

	// publish the movie
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

	// return published movie
	c.JSON(http.StatusCreated, movie)
}
