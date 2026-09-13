package models

import "time"

type Payment struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	BookingID     uint    `json:"booking_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `gorm:"default:'PENDING'" json:"status"`

	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpayPaymentID string `json:"razorpay_payment_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
