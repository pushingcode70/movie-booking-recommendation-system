package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

// registerUserGenreRoutes registers routes for managing the authenticated user's favorite movie genres
func RegisterUserGenreRoutes(
	router *gin.Engine,
	handler *handlers.UserGenreHandler,
) {

	userGenres := router.Group("/users/me/genres")
	userGenres.Use(middleware.AuthMiddleware())
	{
		// get all favorite genres
		userGenres.GET("", handler.GetFavoriteGenres)

		// add a genre to favorite genres
		userGenres.POST("/:genreID", handler.AddFavoriteGenre)

		// remove a genre from favorite genres
		userGenres.DELETE("/:genreID", handler.RemoveFavoriteGenre)
	}
}
