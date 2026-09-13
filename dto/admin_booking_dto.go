package dto

import "time"

// adminRecentBookingResponse represents a booking displayed in the admin dashboard
type AdminRecentBookingResponse struct {
	BookingID     uint      `json:"booking_id"`
	CustomerName  string    `json:"customer_name"`
	MovieTitle    string    `json:"movie_title"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	ShowID        uint      `json:"show_id"`
	ShowStartTime time.Time `json:"show_start_time"`
	ShowEndTime   time.Time `json:"show_end_time"`
}
