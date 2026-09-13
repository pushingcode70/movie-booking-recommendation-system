package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSeatRoutes(router *gin.Engine, seatHandler *handlers.SeatHandler) {

	// Public routes
	router.GET("/seats", seatHandler.GetAllSeats)
	router.GET("/seats/:id", seatHandler.GetSeatByID)
	router.GET("/shows/:id/seats", seatHandler.GetShowSeatLayout)

	// Admin-only route
	admin := router.Group("/")
	admin.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)

	admin.POST("/screens/:id/seats/generate", seatHandler.GenerateSeats)

	// Existing seat CRUD
	seats := router.Group("/seats")
	seats.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		seats.POST("", seatHandler.CreateSeat)
		seats.PUT("/:id", seatHandler.UpdateSeat)
		seats.DELETE("/:id", seatHandler.DeleteSeat)
	}
}
