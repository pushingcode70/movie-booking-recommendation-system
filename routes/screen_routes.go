package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterScreenRoutes(router *gin.Engine, screenHandler *handlers.ScreenHandler) {

	// public routes
	router.GET("/screens", screenHandler.GetAllScreens)
	router.GET("/screens/:id", screenHandler.GetScreenByID)

	// protected routes
	screens := router.Group("/screens")
	screens.Use(middleware.AuthMiddleware(),
		middleware.AdminMiddleware())
	{
		screens.POST("", screenHandler.CreateScreen)
		screens.PUT("/:id", screenHandler.UpdateScreen)
		screens.DELETE("/:id", screenHandler.DeleteScreen)
	}
}
