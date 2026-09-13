package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPaymentRoutes(router *gin.Engine, handler *handlers.PaymentHandler) {

	payments := router.Group("/payments")
	payments.Use(middleware.AuthMiddleware())
	{

		payments.POST("/verify", handler.VerifyPayment)
	}
}
