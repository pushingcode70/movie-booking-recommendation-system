package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(router *gin.Engine, bookingHandler *handlers.BookingHandler) {

	bookings := router.Group("/bookings")
	bookings.Use(middleware.AuthMiddleware())
	{
		bookings.POST("", bookingHandler.CreateBooking)

		bookings.GET("/:id", bookingHandler.GetBookingByID)
		bookings.GET("/user/:user_id", bookingHandler.GetBookingsByUserID)
	}

	adminBookings := router.Group("/bookings")
	adminBookings.Use(
		middleware.AuthMiddleware(),
		middleware.AdminMiddleware(),
	)
	{
		adminBookings.PUT("/:id/confirm", bookingHandler.ConfirmBooking)
		adminBookings.DELETE("/:id", bookingHandler.CancelBooking)
	}
}
