package models

import "time"

type WatchedMovie struct {
	ID uint `gorm:"primaryKey"`

	UserID uint `gorm:"not null;uniqueIndex:idx_user_tmdb_watched"`

	TMDBID int `gorm:"not null;uniqueIndex:idx_user_tmdb_watched"`

	Rating *int

	Review *string `gorm:"type:text"`

	WatchedAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
