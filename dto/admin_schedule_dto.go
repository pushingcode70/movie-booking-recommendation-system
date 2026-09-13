package dto

import "time"

// Represents one scheduled show in a theatre.
type AdminScheduleResponse struct {
	ShowID       uint      `json:"show_id"`
	MovieID      uint      `json:"movie_id"`
	MovieTitle   string    `json:"movie_title"`
	ScreenID     uint      `json:"screen_id"`
	ScreenNumber int       `json:"screen_number"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Price        float64   `json:"price"`
}
