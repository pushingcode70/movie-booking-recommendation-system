package routes

import (
	"movie-booking/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterGenreRoutes registers all genre routes.
func RegisterGenreRoutes(router *gin.Engine, handler *handlers.GenreHandler) {

	genres := router.Group("/genres")
	{
		// Get all available movie genres.
		genres.GET("", handler.GetAllGenres)
	}
}
