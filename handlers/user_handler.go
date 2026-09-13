package handlers

import (
	"net/http"

	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

// constructor
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// getmyprofile of currently logged in user
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	//get authenticated user's id from jwt claims stores in gin's context.
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not authenticated",
		})
		return
	}

	//call service to fetch the user's details
	// c.Get() returns an interface{}.
	// The value is type-asserted to its actual type and converted to uint,
	// because GetUserByID expects a uint ID.
	user, err := h.service.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateMyProfile updates the currently logged-in user's profile.
func (h *UserHandler) UpdateMyProfile(c *gin.Context) {

	// Get the authenticated user's ID from the Gin context.
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	var req dto.UpdateUserRequest

	// Bind the updated data from the request body.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user := models.User{
		ID:    userID.(uint),
		Name:  req.Name,
		Email: req.Email,
	}

	// Always use the authenticated user's ID.
	// Never trust an ID sent by the client.

	// Update the user.
	if err := h.service.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
