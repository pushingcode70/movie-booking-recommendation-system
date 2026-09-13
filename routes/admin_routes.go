package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(
	router *gin.Engine,
	handler *handlers.AdminHandler,
) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())

	{
		admin.GET("/dashboard", handler.GetDashboardStats)
		admin.GET("/payments/recent", handler.GetRecentPayments)
		admin.GET("/bookings/recent", handler.GetRecentBookings)
		admin.GET("/shows/running", handler.GetRunningShows)
		admin.GET("/theatres/:id/schedule", handler.GetTheatreSchedule)
	}
}
