package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
	"movie-booking/utils"
)

type AuthService struct {
	repo         *repositories.UserRepository
	emailService *EmailService
}

// Constructor
func NewAuthService(repo *repositories.UserRepository, emailService *EmailService) *AuthService {
	return &AuthService{
		repo:         repo,
		emailService: emailService,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password.
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Prevent unverified users from logging in.
	if !user.IsVerified {
		return "", errors.New("please verify your email before logging in")
	}

	// Generate JWT.
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(req dto.RegisterRequest) error {

	// Check if the email is already registered.
	existingUser, err := s.repo.GetUserByEmail(req.Email)

	if err == nil && existingUser != nil {
		if !existingUser.IsVerified {
			// Account exists but is not verified yet.
			// Update credentials, generate fresh OTP, and send verification email.
			hashedPassword, err := utils.HashPassword(req.Password)
			if err != nil {
				return err
			}

			otp := fmt.Sprintf("%06d", rand.Intn(1000000))
			expiry := time.Now().Add(10 * time.Minute)

			existingUser.Name = req.Name
			existingUser.Password = hashedPassword
			existingUser.OTPCode = otp
			existingUser.OTPExpiresAt = &expiry

			err = s.repo.UpdateUser(existingUser)
			if err != nil {
				return err
			}

			return s.emailService.SendVerificationOTP(dto.OTPEmailData{
				Title:   "Verify Your Email",
				Message: "Use the OTP below to verify your account.",
				OTP:     otp,
				ToEmail: existingUser.Email,
			})
		}
		return errors.New("email already registered")
	}

	//  Hash password before storing it
	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	//  Generate OTP
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	//otp expired after ten minutes
	expiry := time.Now().Add(10 * time.Minute)

	//  Create user
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         models.RoleCustomer,
		IsVerified:   false,
		OTPCode:      otp,
		OTPExpiresAt: &expiry,
	}

	//  Save user
	err = s.repo.CreateUser(&user)
	if err != nil {
		return err
	}
	// Send verification email.
	err = s.emailService.SendVerificationOTP(dto.OTPEmailData{
		Title:   "Verify Your Email",
		Message: "Use the OTP below to verify your account.",
		OTP:     otp,
		ToEmail: user.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(req dto.VerifyOTPRequest) error {

	// Find the user by email.
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if the email is already verified.
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Verify the OTP.
	if user.OTPCode != req.OTP {
		return errors.New("invalid OTP")
	}

	// Check if an OTP exists and has not expired.
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return errors.New("OTP has expired")
	}

	// Mark the email as verified.
	user.IsVerified = true

	// Clear OTP after successful verification.
	user.OTPCode = ""
	user.OTPExpiresAt = nil

	// Save the updated user.
	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

// resend otp
func (s *AuthService) ResendOTP(req dto.ResendOTPRequest) error {

	// Find the user by email.
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if the email is already verified.
	if user.IsVerified {
		return errors.New("email is already verified")
	}

	// Generate a new 6-digit OTP.
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// OTP expires after 10 minutes.
	expiry := time.Now().Add(10 * time.Minute)

	// Update the user's OTP.
	user.OTPCode = otp
	user.OTPExpiresAt = &expiry

	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	// Send the verification email.
	err = s.emailService.SendVerificationOTP(dto.OTPEmailData{
		Title:   "Verify Your Email",
		Message: "Use the OTP below to verify your account.",
		OTP:     otp,
		ToEmail: user.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

// ForgotPassword starts the password reset process.
//
// Flow:
// 1. User clicks "Forgot Password?" on the login page.
// 2. User enters their email address.
// 3. Backend verifies that the account exists.
// 4. Backend generates and stores a temporary OTP.
// 5. Backend emails the OTP to the user.
//
// The password is NOT changed here.
// The actual password update happens later in ResetPassword(),
// after the user submits the OTP and a new password.

func (s *AuthService) ForgotPassword(req dto.ForgotPasswordRequest) error {

	// Find the user by email.
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate a new 6-digit OTP.
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// OTP expires after 10 minutes.
	expiry := time.Now().Add(10 * time.Minute)

	// Save the new OTP.
	user.OTPCode = otp
	user.OTPExpiresAt = &expiry

	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	// Send password reset OTP email.
	err = s.emailService.SendVerificationOTP(dto.OTPEmailData{
		Title:   "Reset Your Password",
		Message: "Use the OTP below to reset your password.",
		OTP:     otp,
		ToEmail: user.Email,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResetPassword(req dto.ResetPasswordRequest) error {

	// Find the user by email.
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify the OTP.
	if user.OTPCode != req.OTP {
		return errors.New("invalid OTP")
	}

	// Check if the OTP has expired.
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return errors.New("OTP has expired")
	}

	// Hash the new password.
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// Update the password.
	user.Password = hashedPassword

	// Clear the OTP.
	user.OTPCode = ""
	user.OTPExpiresAt = nil

	// Save the updated user.
	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
