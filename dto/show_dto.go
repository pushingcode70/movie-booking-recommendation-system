package dto

import "time"

type CreateShowRequest struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	ScreenID  uint      `json:"screen_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	Price     float64   `json:"price" binding:"required,gt=0"`
}

type UpdateShowRequest struct {
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}

type ShowResponse struct {
	ID        uint      `json:"id"`
	MovieID   uint      `json:"movie_id"`
	ScreenID  uint      `json:"screen_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Price     float64   `json:"price"`
}
