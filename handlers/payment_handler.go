package handlers

import (
	"movie-booking/services"
	"net/http"

	"github.com/gin-gonic/gin"

	"movie-booking/dto"
)

type PaymentHandler struct {
	service *services.PaymentService
}

func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: service,
	}
}

func (h *PaymentHandler) VerifyPayment(c *gin.Context) {
	var req dto.VerifyPaymentRequest

	// read json request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// verify payment
	err := h.service.VerifyPayment(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment verified successfully",
	})

}
