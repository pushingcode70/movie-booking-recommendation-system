package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, movieHandler *handlers.MovieHandler) {

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Movie Booking API is running",
		})
	})

	// Public
	router.GET("/movies", movieHandler.GetAllMovies)
	router.GET("/movies/tmdb/:tmdbId", movieHandler.GetMovieByTMDBID)
	router.GET("/movies/:id", movieHandler.GetMovieByID)

	// Protected
	movies := router.Group("/movies")
	movies.Use(middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		movies.POST("", movieHandler.CreateMovie)
		movies.PUT("/:id", movieHandler.UpdateMovie)
		movies.DELETE("/:id", movieHandler.DeleteMovie)
	}
}
