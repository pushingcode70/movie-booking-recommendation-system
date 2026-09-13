package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// registerEmbeddingRoutes registers routes for movie embeddings
func RegisterEmbeddingRoutes(
	router *gin.Engine,
	handler *handlers.EmbeddingHandler,
) {

	embeddings := router.Group("/embeddings")
	embeddings.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// generate and store an embedding for a tmdb movie
		embeddings.POST("/movie/:tmdbID", handler.GenerateMovieEmbedding)
	}
}
