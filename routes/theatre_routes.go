package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTheatreRoutes(router *gin.Engine, theatreHandler *handlers.TheatreHandler) {

	// Public routes
	router.GET("/theatres", theatreHandler.GetAllTheatres)
	router.GET("/theatres/:id", theatreHandler.GetTheatreByID)
	router.GET("/theatres/:id/schedule", theatreHandler.GetTheatreSchedule)

	// Protected routes
	theatres := router.Group("/theatres")
	theatres.Use(middleware.AuthMiddleware(),
		middleware.AdminMiddleware())
	{
		theatres.POST("", theatreHandler.CreateTheatre)
		theatres.PUT("/:id", theatreHandler.UpdateTheatre)
		theatres.DELETE("/:id", theatreHandler.DeleteTheatre)
	}
}
