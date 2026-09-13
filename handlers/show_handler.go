package handlers

import (
	"net/http"
	"strconv"
	"time"

	"movie-booking/models"
	"movie-booking/services"

	"movie-booking/dto"

	"github.com/gin-gonic/gin"
)

type ShowHandler struct {
	service *services.ShowService
}

func NewShowHandler(service *services.ShowService) *ShowHandler {
	return &ShowHandler{service: service}
}

func (h *ShowHandler) CreateShow(c *gin.Context) {

	var req dto.CreateShowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	show := models.Show{
		MovieID:   req.MovieID,
		ScreenID:  req.ScreenID,
		StartTime: req.StartTime,
		Price:     req.Price,
	}

	if err := h.service.CreateShow(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	c.JSON(http.StatusCreated, show)
}

func (h *ShowHandler) GetShowByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	show, err := h.service.GetShowByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	// block access if the show has already ended
	if show.EndTime.Before(time.Now()) {
		c.JSON(http.StatusGone, gin.H{"error": "This show has already ended and is no longer available for booking"})
		return
	}

	c.JSON(http.StatusOK, show)
}

func (h *ShowHandler) GetAllShows(c *gin.Context) {
	shows, err := h.service.GetAllShows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func (h *ShowHandler) UpdateShow(c *gin.Context) {
	// parsing and binding received request
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req dto.UpdateShowRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get existing show
	show, err := h.service.GetShowByID(uint(id))

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	// update editable fields
	show.StartTime = req.StartTime
	show.Price = req.Price

	// save changes
	if err := h.service.UpdateShow(show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, show)
}

func (h *ShowHandler) DeleteShow(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteShow(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Show deleted successfully"})
}

func (h *ShowHandler) GetShowHistoryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	show, err := h.service.GetShowByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Show not found"})
		return
	}

	c.JSON(http.StatusOK, show)
}
