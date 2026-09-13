package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/models"
	"movie-booking/services"

	"movie-booking/dto"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	service *services.SeatService
}

// Constructor
func NewSeatHandler(service *services.SeatService) *SeatHandler {
	return &SeatHandler{service: service}
}

// Create Seat
func (h *SeatHandler) CreateSeat(c *gin.Context) {
	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, seat)
}

// Get Seat by ID
func (h *SeatHandler) GetSeatByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	seat, err := h.service.GetSeatByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Seat not found"})
		return
	}

	c.JSON(http.StatusOK, seat)
}

// Get All Seats
func (h *SeatHandler) GetAllSeats(c *gin.Context) {
	seats, err := h.service.GetAllSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}

// Update Seat
func (h *SeatHandler) UpdateSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var seat models.Seat

	if err := c.ShouldBindJSON(&seat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	seat.ID = uint(id)

	if err := h.service.UpdateSeat(&seat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seat)
}

// Delete Seat
func (h *SeatHandler) DeleteSeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteSeat(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seat deleted successfully"})
}

// GetShowSeatLayout returns the seat layout for a show.
func (h *SeatHandler) GetShowSeatLayout(c *gin.Context) {

	// Read show ID from URL.
	showID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid show ID",
		})
		return
	}

	// Call service.
	layout, err := h.service.GetShowSeatLayout(uint(showID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// No show or no seats found.
	if layout == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Show not found",
		})
		return
	}

	c.JSON(http.StatusOK, layout)
}

// GenerateSeats creates a seat layout for a screen.
func (h *SeatHandler) GenerateSeats(c *gin.Context) {

	// Get screen ID from URL.
	screenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid screen ID",
		})
		return
	}

	// Read request body.
	var req dto.GenerateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Generate the seat layout.
	err = h.service.GenerateSeats(uint(screenID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Seat layout generated successfully",
		"total_seats": req.Rows * req.SeatsPerRow,
	})
}
