package models

import "time"

const (
	RoleCustomer = "customer"
	RoleAdmin    = "admin"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"-"`
	Role     string `gorm:"type:varchar(20);default:'customer'" json:"role"`

	// Genres selected by the user as their preferred movie genres.
	FavoriteGenres []Genre `gorm:"many2many:user_favorite_genres;" json:"favorite_genres"`

	// Indicates whether the user's email has been verified.
	IsVerified bool `gorm:"default:false" json:"is_verified"`

	// Stores the latest OTP sent to the user's email.
	// Cleared after successful verification.
	OTPCode string `json:"-"`

	// Expiration time for the OTP.
	// Nil means there is no active OTP.
	OTPExpiresAt *time.Time `json:"-"`
}
