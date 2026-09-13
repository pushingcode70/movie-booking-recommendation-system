package handlers

import (
	"movie-booking/dto"
	"movie-booking/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	service *services.BookingService
}

// Constructor
func NewBookingHandler(service *services.BookingService) *BookingHandler {
	return &BookingHandler{
		service: service,
	}
}

// Create Booking
func (h *BookingHandler) CreateBooking(c *gin.Context) {

	// DTO to receive only the fields the client is allowed to send.
	// We don't bind directly to the Booking model for security reasons.
	var req dto.CreateBookingRequest

	// Read and validate JSON request body.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get the authenticated user's ID from Gin's context.
	// The AuthMiddleware stored it there after validating the JWT.
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	// JWT MapClaims stores numeric values as float64.
	// Convert it back to uint before passing it to the service.
	userID := userIDValue.(uint)
	// Call the service, which contains all booking business logic.
	response, err := h.service.BookService(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Return the created booking.
	c.JSON(http.StatusCreated, response) //response is of type *dto.BookingResponse
}

func (h *BookingHandler) GetBookingByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userID := userIDValue.(uint)

	booking, err := h.service.GetBookingByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	if booking.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not allowed to access this booking",
		})
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingHandler) GetBookingsByUserID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	userID := userIDValue.(uint)

	if uint(id) != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not allowed to access these bookings",
		})
		return
	}

	bookings, err := h.service.GetBookingsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

// confirm Booking
func (h *BookingHandler) ConfirmBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := h.service.ConfirmBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking Confirmed",
	})
}

// cancel Booking
func (h *BookingHandler) CancelBooking(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	if err := h.service.CancelBooking(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Booking Cancelled",
	})
}
