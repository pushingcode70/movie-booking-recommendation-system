package handlers

import (
	"net/http"
	"strconv"

	"movie-booking/models"
	"movie-booking/services"

	"github.com/gin-gonic/gin"
)

type ScreenHandler struct {
	service *services.ScreenService
}

// Constructor
func NewScreenHandler(service *services.ScreenService) *ScreenHandler {
	return &ScreenHandler{service: service}
}

// Create Screen
func (h *ScreenHandler) CreateScreen(c *gin.Context) {
	var screen models.Screen

	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateScreen(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, screen)
}

// Get Screen by ID
func (h *ScreenHandler) GetScreenByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	screen, err := h.service.GetScreenByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Screen not found"})
		return
	}

	c.JSON(http.StatusOK, screen)
}

// Get All Screens
func (h *ScreenHandler) GetAllScreens(c *gin.Context) {
	screens, err := h.service.GetAllScreens()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, screens)
}

// Update Screen
func (h *ScreenHandler) UpdateScreen(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	var screen models.Screen

	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	screen.ID = uint(id)

	if err := h.service.UpdateScreen(&screen); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, screen)
}

// Delete Screen
func (h *ScreenHandler) DeleteScreen(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid screen ID"})
		return
	}

	if err := h.service.DeleteScreen(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Screen deleted successfully",
	})
}
