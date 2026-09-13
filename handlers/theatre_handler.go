package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/models"
	"movie-booking/services"

	"time"

	"github.com/gin-gonic/gin"
)

type TheatreHandler struct {
	service *services.TheatreService
}

// Constructor
func NewTheatreHandler(service *services.TheatreService) *TheatreHandler {
	return &TheatreHandler{service: service}
}

// Create Theatre
func (h *TheatreHandler) CreateTheatre(c *gin.Context) {
	var theatre models.Theatre

	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTheatre(&theatre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, theatre)
}

// Get Theatre by ID
func (h *TheatreHandler) GetTheatreByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	theatre, err := h.service.GetTheatreByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theatre not found"})
		return
	}

	c.JSON(http.StatusOK, theatre)
}

func (h *TheatreHandler) GetAllTheatres(c *gin.Context) {

	theatres, err := h.service.GetAllTheatres()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch theatres",
		})
		return
	}

	c.JSON(http.StatusOK, theatres)
}

// Update Theatre
func (h *TheatreHandler) UpdateTheatre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	var theatre models.Theatre

	if err := c.ShouldBindJSON(&theatre); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	theatre.ID = uint(id)

	if err := h.service.UpdateTheatre(&theatre); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, theatre)
}

// Delete Theatre
func (h *TheatreHandler) DeleteTheatre(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theatre ID"})
		return
	}

	if err := h.service.DeleteTheatre(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Theatre deleted successfully",
	})
}

// GetTheatreSchedule returns the customer-facing schedule
// for a theatre on a given date.
func (h *TheatreHandler) GetTheatreSchedule(c *gin.Context) {

	// Read theatre id from URL.
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid theatre id",
		})
		return
	}

	// Read date from query parameter.
	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date is required",
		})
		return
	}

	// Validate date format.
	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format (YYYY-MM-DD)",
		})
		return
	}

	// Fetch theatre schedule.
	schedule, err := h.service.GetTheatreSchedule(uint(id), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch theatre schedule",
		})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
