package dto

import "time"

// adminRunningShowResponse represents a currently running show displayed on the admin dashboard
type AdminRunningShowResponse struct {
	ShowID       uint      `json:"show_id"`
	MovieTitle   string    `json:"movie_title"`
	TheatreName  string    `json:"theatre_name"`
	ScreenNumber int       `json:"screen_number"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`

	BookedSeats int64 `json:"booked_seats"`
	TotalSeats  int64 `json:"total_seats"`

	Occupancy float64 `json:"occupancy"`
}
