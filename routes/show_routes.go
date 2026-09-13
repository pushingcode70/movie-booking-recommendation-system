package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterShowRoutes(router *gin.Engine, showHandler *handlers.ShowHandler) {

	// public routes
	router.GET("/shows", showHandler.GetAllShows)
	router.GET("/shows/:id", showHandler.GetShowByID)
	router.GET("/shows/history/:id", showHandler.GetShowHistoryByID)

	// protected routes
	shows := router.Group("/shows")
	shows.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		shows.POST("", showHandler.CreateShow)
		shows.PUT("/:id", showHandler.UpdateShow)
		shows.DELETE("/:id", showHandler.DeleteShow)
	}
}
