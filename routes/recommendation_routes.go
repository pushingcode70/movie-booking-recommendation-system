package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// registerRecommendationRoutes registers routes for recommendation features
func RegisterRecommendationRoutes(
	router *gin.Engine,
	handler *handlers.RecommendationHandler,
) {

	recommendations := router.Group("/recommendations")
	recommendations.Use(middleware.AuthMiddleware())
	{
		// get user's favorite genres
		recommendations.GET("/genres", handler.GetFavoriteGenreRecommendations)

		// get user's wishlist and watched movie ids
		recommendations.GET("/taste", handler.GetTasteContext)

		recommendations.POST("/genre", handler.GetGenreRecommendations)

		// for custom prompt
		recommendations.POST("/custom", handler.GetPromptRecommendations)

		recommendations.POST("/prompt-genre", handler.GetPromptGenreRecommendations)

		recommendations.POST("/movie", handler.GetMovieRecommendations)

	}
}
