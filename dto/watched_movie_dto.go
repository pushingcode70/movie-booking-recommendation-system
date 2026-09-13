package dto

import (
	"movie-booking/models"
	"time"
)

type WatchedMovieItemResponse struct {
	Movie     models.Movie `json:"movie"`
	Rating    *int         `json:"rating"`
	Review    *string      `json:"review"`
	WatchedAt time.Time    `json:"watched_at"`
}
