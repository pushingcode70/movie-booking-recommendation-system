package models

import "time"

type Booking struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `json:"user_id"`
	ShowID      uint    `json:"show_id"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `gorm:"default:'PENDING'" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
