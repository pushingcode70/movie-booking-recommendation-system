package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWatchedMovieRoutes(router *gin.Engine, handler *handlers.WatchedMovieHandler) {
	watched := router.Group("/watched")
	watched.Use(middleware.AuthMiddleware())
	{
		watched.POST("", handler.AddWatchedMovie)
		watched.PUT("/:tmdbID", handler.UpdateWatchedMovie)
		watched.GET("", handler.GetWatchedMovies)
		watched.DELETE("/:tmdbID", handler.RemoveWatchedMovie)
	}
}
