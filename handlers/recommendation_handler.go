package handlers

import (
	"net/http"

	"movie-booking/dto"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type RecommendationHandler struct {
	service *services.RecommendationService
}

type CustomRecommendationRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

func NewRecommendationHandler(service *services.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{
		service: service,
	}
}

func (h *RecommendationHandler) GetFavoriteGenreRecommendations(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	genres, err := h.service.GetFavoriteGenreContext(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *RecommendationHandler) GetTasteContext(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)

	wishlistMovies, watchedMovies, err := h.service.GetTasteContext(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wishlist_movie_ids": wishlistMovies,
		"watched_movie_ids":  watchedMovies,
	})
}

func (h *RecommendationHandler) GetPromptRecommendations(c *gin.Context) {

	var req CustomRecommendationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	recommendations, err := h.service.GetPromptRecommendations(req.Prompt, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func (h *RecommendationHandler) GetGenreRecommendations(c *gin.Context) {

	var req dto.RecommendationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// ask service to generate balanced set of recommendations for selected genres
	recommendations, err := h.service.GetGenreRecommendations(
		req.GenreIDs,
		30,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func (h *RecommendationHandler) GetPromptGenreRecommendations(c *gin.Context) {

	var req dto.RecommendationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// prompt + genre requires both fields
	if req.Prompt == "" || len(req.GenreIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "prompt and genre_ids are required",
		})
		return
	}

	recommendations, err := h.service.GetPromptGenreRecommendations(req.Prompt, req.GenreIDs, 30)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}

func (h *RecommendationHandler) GetMovieRecommendations(c *gin.Context) {

	var req dto.RecommendationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	if len(req.MovieIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "movie_ids are required",
		})
		return
	}

	if len(req.GenreIDs) > 0 || req.Prompt != "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "movie recommendations cannot contain genres or prompt",
		})
		return
	}

	recommendations, err := h.service.GetMovieRecommendations(req.MovieIDs, 30)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, recommendations)
}
