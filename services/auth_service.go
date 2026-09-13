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

	// verify password
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// prevent unverified users from logging in
	if !user.IsVerified {
		return "", errors.New("please verify your email before logging in")
	}

	// generate jwt
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Register(req dto.RegisterRequest) error {

	// check if the email is already registered
	existingUser, err := s.repo.GetUserByEmail(req.Email)

	if err == nil && existingUser != nil {
		if !existingUser.IsVerified {
			// account exists but is not verified yet
			// update credentials, generate fresh otp, and send verification email
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

	// hash password before storing it
	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		return err
	}

	// generate otp
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// otp expires after ten minutes
	expiry := time.Now().Add(10 * time.Minute)

	// create user
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPassword,
		Role:         models.RoleCustomer,
		IsVerified:   false,
		OTPCode:      otp,
		OTPExpiresAt: &expiry,
	}

	// save user
	err = s.repo.CreateUser(&user)
	if err != nil {
		return err
	}
	// send verification email
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

	// find the user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// check if the email is already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// verify the otp
	if user.OTPCode != req.OTP {
		return errors.New("invalid OTP")
	}

	// check if an otp exists and has not expired
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return errors.New("OTP has expired")
	}

	// mark the email as verified
	user.IsVerified = true

	// clear otp after successful verification
	user.OTPCode = ""
	user.OTPExpiresAt = nil

	// save the updated user
	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ResendOTP(req dto.ResendOTPRequest) error {

	// find the user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// check if the email is already verified
	if user.IsVerified {
		return errors.New("email is already verified")
	}

	// generate a new 6-digit otp
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// otp expires after 10 minutes
	expiry := time.Now().Add(10 * time.Minute)

	// update the user's otp
	user.OTPCode = otp
	user.OTPExpiresAt = &expiry

	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	// send the verification email
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

func (s *AuthService) ForgotPassword(req dto.ForgotPasswordRequest) error {

	// find the user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// generate a new 6-digit otp
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	// otp expires after 10 minutes
	expiry := time.Now().Add(10 * time.Minute)

	// save the new otp
	user.OTPCode = otp
	user.OTPExpiresAt = &expiry

	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	// send password reset otp email
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

	// find the user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// verify the otp
	if user.OTPCode != req.OTP {
		return errors.New("invalid OTP")
	}

	// check if the otp has expired
	if user.OTPExpiresAt == nil || time.Now().After(*user.OTPExpiresAt) {
		return errors.New("OTP has expired")
	}

	// hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	// update the password
	user.Password = hashedPassword

	// clear the otp
	user.OTPCode = ""
	user.OTPExpiresAt = nil

	// save the updated user
	err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
