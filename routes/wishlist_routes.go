package routes

import (
	"movie-booking/handlers"
	"movie-booking/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWishlistRoutes(router *gin.Engine, handler *handlers.WishlistHandler) {
	wishlist := router.Group("/wishlist")
	wishlist.Use(middleware.AuthMiddleware())
	{
		wishlist.POST("", handler.AddMovie)
		wishlist.GET("", handler.GetWishlist)
		wishlist.DELETE("/:tmdbID", handler.RemoveMovie)
	}
}
