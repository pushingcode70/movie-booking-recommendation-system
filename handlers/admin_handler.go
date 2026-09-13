package handlers

import (
	"net/http"

	"movie-booking/services"

	"github.com/gin-gonic/gin"

	"strconv"
)

type AdminHandler struct {
	service *services.AdminService
}

func NewAdminHandler(service *services.AdminService) *AdminHandler {
	return &AdminHandler{
		service: service,
	}
}

func (h *AdminHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *AdminHandler) GetRecentPayments(c *gin.Context) {
	date := c.Query("date") // optional, e.g. ?date=2026-07-26

	payments, err := h.service.GetRecentPayments(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch payments",
		})
		return
	}

	c.JSON(http.StatusOK, payments)
}

func (h *AdminHandler) GetRecentBookings(c *gin.Context) {
	date := c.Query("date")

	bookings, err := h.service.GetRecentBookings(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recent bookings",
		})
		return
	}

	c.JSON(http.StatusOK, bookings)
}

func (h *AdminHandler) GetRunningShows(c *gin.Context) {

	shows, err := h.service.GetRunningShows()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch running shows",
		})
		return
	}

	c.JSON(http.StatusOK, shows)
}

func (h *AdminHandler) GetTheatreSchedule(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid theatre id",
		})
		return
	}

	date := c.Query("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date query parameter is required (YYYY-MM-DD)",
		})
		return
	}

	schedule, err := h.service.GetTheatreSchedule(uint(id), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch theatre schedule",
		})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
