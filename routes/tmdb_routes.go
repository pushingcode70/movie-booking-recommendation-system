package routes

import (
	"movie-booking/handlers"

	"github.com/gin-gonic/gin"
)

// registerTMDBRoutes registers all tmdb routes
func RegisterTMDBRoutes(router *gin.Engine, handler *handlers.TMDBHandler, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) {

	tmdb := router.Group("/tmdb")
	{

		// search movies by title
		tmdb.GET("/search", handler.SearchMovies)

		// search by tmdb id
		tmdb.GET("/movie/:id", handler.GetMovieDetails)

		// admin-only route to publish a tmdb movie
		tmdb.POST(
			"/publish/:id",
			authMiddleware,
			adminMiddleware,
			handler.PublishMovie,
		)
	}
}
