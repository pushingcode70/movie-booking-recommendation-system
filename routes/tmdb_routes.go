package routes

import (
	"movie-booking/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterTMDBRoutes registers all TMDb routes.
func RegisterTMDBRoutes(router *gin.Engine, handler *handlers.TMDBHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {

	tmdb := router.Group("/tmdb")
	{

		//search movies by title.
		tmdb.GET("/search", handler.SearchMovies)

		//by id
		tmdb.GET("/movie/:id", handler.GetMovieDetails)

		// Admin-only route to publish a TMDb movie
		tmdb.POST(
			"/publish/:id",
			authMiddleware,
			adminMiddleware,
			handler.PublishMovie,
		)
	}
}
