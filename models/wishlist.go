package models

import "time"

type Wishlist struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;uniqueIndex:idx_user_tmdb"`

	TMDBID int `gorm:"not null;uniqueIndex:idx_user_tmdb"`

	CreatedAt time.Time
}
