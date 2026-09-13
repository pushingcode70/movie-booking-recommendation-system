package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/dto"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type WishlistHandler struct {
	service *services.WishlistService
}

func NewWishlistHandler(service *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		service: service,
	}
}

// Add a movie to the wishlist.
func (h *WishlistHandler) AddMovie(c *gin.Context) {

	var req dto.CreateWishlistRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	if err := h.service.AddMovie(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Movie added to wishlist successfully",
	})
}

// Get the user's wishlist.
func (h *WishlistHandler) GetWishlist(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)
	response, err := h.service.GetWishlist(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// / Remove a movie from the wishlist.
func (h *WishlistHandler) RemoveMovie(c *gin.Context) {

	userID := c.MustGet("user_id").(uint)
	tmdbID, err := strconv.Atoi(c.Param("tmdbID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid TMDb ID",
		})
		return
	}

	if err := h.service.RemoveMovie(userID, tmdbID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie removed from wishlist successfully",
	})
}
