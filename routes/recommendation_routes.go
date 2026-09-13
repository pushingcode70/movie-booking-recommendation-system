package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRecommendationRoutes registers routes for recommendation features.
func RegisterRecommendationRoutes(
	router *gin.Engine,
	handler *handlers.RecommendationHandler,
) {

	recommendations := router.Group("/recommendations")
	recommendations.Use(middleware.AuthMiddleware())
	{
		// get the user's favorite genres.
		recommendations.GET("/genres", handler.GetFavoriteGenreRecommendations)

		// get the user's wishlist and watched movie IDs.
		recommendations.GET("/taste", handler.GetTasteContext)

		recommendations.POST("/genre", handler.GetGenreRecommendations)

		//for custom promt
		recommendations.POST("/custom", handler.GetPromptRecommendations)

		recommendations.POST("/prompt-genre", handler.GetPromptGenreRecommendations)

		recommendations.POST("/movie", handler.GetMovieRecommendations)

	}
}
