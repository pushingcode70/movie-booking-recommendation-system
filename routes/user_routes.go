package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userHandler *handlers.UserHandler) {

	users := router.Group("/users")
	users.Use(middleware.AuthMiddleware())
	{
		//logged in user can view and update their own profile
		users.GET("/me", userHandler.GetMyProfile)
		users.PUT("/me", userHandler.UpdateMyProfile)

	}
}
