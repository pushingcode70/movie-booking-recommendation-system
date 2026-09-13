package dto

import "time"

// AdminRecentPaymentResponse represents a payment shown
// on the admin dashboard's "Recent Payments" widget.
type AdminRecentPaymentResponse struct {
	PaymentID    uint      `json:"payment_id"`
	BookingID    uint      `json:"booking_id"`
	CustomerName string    `json:"customer_name"`
	MovieTitle   string    `json:"movie_title"`
	Amount       float64   `json:"amount"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}
