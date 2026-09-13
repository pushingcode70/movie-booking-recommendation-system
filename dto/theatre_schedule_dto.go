package dto

import "time"

type TheatreScheduleResponse struct {
	TheatreID   uint                   `json:"theatre_id"`
	TheatreName string                 `json:"theatre_name"`
	Address     string                 `json:"address"`
	Date        string                 `json:"date"`
	Movies      []TheatreMovieSchedule `json:"movies"`
}

type TheatreMovieSchedule struct {
	MovieID    uint              `json:"movie_id"`
	Title      string            `json:"title"`
	PosterPath string            `json:"poster_path"`
	Language   string            `json:"language"`
	Duration   int               `json:"duration"`
	ShowTimes  []TheatreShowTime `json:"show_times"`
}

type TheatreShowTime struct {
	ShowID    uint    `json:"show_id"`
	StartTime string  `json:"start_time"`
	Price     float64 `json:"price"`
}

type TheatreScheduleRow struct {
	TheatreID   uint   `json:"theatre_id"`
	TheatreName string `json:"theatre_name"`
	Address     string `json:"address"`

	MovieID    uint   `json:"movie_id"`
	Title      string `json:"title"`
	PosterPath string `json:"poster_path"`
	Language   string `json:"language"`
	Duration   int    `json:"duration"`

	ShowID    uint      `json:"show_id"`
	StartTime time.Time `json:"start_time"`
	Price     float64   `json:"price"`
}
