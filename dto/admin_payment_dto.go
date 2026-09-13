package dto

import "time"

// adminRecentPaymentResponse represents a payment shown on admin dashboard
type AdminRecentPaymentResponse struct {
	PaymentID    uint      `json:"payment_id"`
	BookingID    uint      `json:"booking_id"`
	CustomerName string    `json:"customer_name"`
	MovieTitle   string    `json:"movie_title"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
