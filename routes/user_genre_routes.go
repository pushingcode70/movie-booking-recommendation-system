package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterUserGenreRoutes registers routes for managing
// the authenticated user's favorite movie genres.
func RegisterUserGenreRoutes(
	router *gin.Engine,
	handler *handlers.UserGenreHandler,
) {

	userGenres := router.Group("/users/me/genres")
	userGenres.Use(middleware.AuthMiddleware())
	{
		// Get all favorite genres.
		userGenres.GET("", handler.GetFavoriteGenres)

		// Add a genre to favorite genres.
		userGenres.POST("/:genreID", handler.AddFavoriteGenre)

		// Remove a genre from favorite genres.
		userGenres.DELETE("/:genreID", handler.RemoveFavoriteGenre)
	}
}
