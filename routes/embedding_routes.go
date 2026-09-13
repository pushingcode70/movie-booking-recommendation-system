package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterEmbeddingRoutes registers routes for movie embeddings.
func RegisterEmbeddingRoutes(
	router *gin.Engine,
	handler *handlers.EmbeddingHandler,
) {

	embeddings := router.Group("/embeddings")
	embeddings.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Generate and store an embedding for a TMDB movie.
		embeddings.POST("/movie/:tmdbID", handler.GenerateMovieEmbedding)
	}
}
