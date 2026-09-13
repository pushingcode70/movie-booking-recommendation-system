package dto

type DashboardResponse struct {
	TotalMovies        int64   `json:"total_movies"`
	TotalTheatres      int64   `json:"total_theatres"`
	TotalScreens       int64   `json:"total_screens"`
	TotalShows         int64   `json:"total_shows"`
	TotalBookings      int64   `json:"total_bookings"`
	TotalRevenue       float64 `json:"total_revenue"`
	PendingPayments    int64   `json:"pending_payments"`
	SuccessfulPayments int64   `json:"successful_payments"`
}

type OTPEmailData struct {
	Title   string
	Message string
	OTP     string
	ToEmail string
}
