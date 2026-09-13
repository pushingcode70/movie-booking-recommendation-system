package dto

import (
	"movie-booking/models"
	"time"
)

type WishlistItemResponse struct {
	Movie   models.Movie `json:"movie"`
	AddedAt time.Time    `json:"added_at"`
}
