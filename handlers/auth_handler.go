package handlers

import (
	"github.com/gin-gonic/gin"
	"movie-booking/dto"
	"movie-booking/services"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

// constructor
func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login
func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest

	//read json request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	//call service
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	//login successful
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req dto.RegisterRequest

	// Read and validate JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service.
	err := h.service.Register(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Registration successful.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Please verify your email using the OTP sent to your inbox.",
	})
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {

	var req dto.VerifyOTPRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.VerifyOTP(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email verified successfully.",
	})
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {

	var req dto.ResendOTPRequest

	// Read and validate JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service.
	err := h.service.ResendOTP(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// OTP sent successfully.
	c.JSON(http.StatusOK, gin.H{
		"message": "OTP sent successfully.",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {

	var req dto.ForgotPasswordRequest

	// Read and validate JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service.
	err := h.service.ForgotPassword(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Password reset OTP sent.
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset OTP sent successfully.",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {

	var req dto.ResetPasswordRequest

	// Read and validate JSON request.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Call the service.
	err := h.service.ResetPassword(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Password reset successful.
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully.",
	})
}
