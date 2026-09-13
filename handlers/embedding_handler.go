package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type EmbeddingHandler struct {
	recommendationService *services.RecommendationService
}

// constructor
func NewEmbeddingHandler(recommendationService *services.RecommendationService) *EmbeddingHandler {

	return &EmbeddingHandler{
		recommendationService: recommendationService,
	}
}

// generates and stores an embedding for a tmdb movie
func (h *EmbeddingHandler) GenerateMovieEmbedding(c *gin.Context) {

	tmdbID, err := strconv.Atoi(c.Param("tmdbID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid TMDB movie ID",
		})
		return
	}

	if err := h.recommendationService.GenerateAndStoreMovieEmbedding(tmdbID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "movie embedding generated and stored successfully",
		"tmdb_id": tmdbID,
	})
}
