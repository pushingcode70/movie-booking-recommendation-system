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

func NewSeatHandler(service *services.SeatService) *SeatHandler {
	return &SeatHandler{service: service}
}

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

func (h *SeatHandler) GetAllSeats(c *gin.Context) {
	seats, err := h.service.GetAllSeats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}

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

func (h *SeatHandler) GetShowSeatLayout(c *gin.Context) {

	// read show ID from URL
	showID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid show ID",
		})
		return
	}

	// call service
	layout, err := h.service.GetShowSeatLayout(uint(showID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// no show or no seats found
	if layout == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Show not found",
		})
		return
	}

	c.JSON(http.StatusOK, layout)
}

func (h *SeatHandler) GenerateSeats(c *gin.Context) {

	// get screen ID from URL
	screenID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid screen ID",
		})
		return
	}

	// read request body
	var req dto.GenerateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// generate seat layout
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
