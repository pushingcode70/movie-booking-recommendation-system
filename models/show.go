package models

import (
	"time"
)

type Show struct {
	ID uint `gorm:"primaryKey" json:"id"`

	MovieID uint  `json:"movie_id"`
	Movie   Movie `gorm:"foreignKey:MovieID" json:"movie"`

	ScreenID uint   `json:"screen_id"`
	Screen   Screen `gorm:"foreignKey:ScreenID" json:"screen"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	Price float64 `json:"price"`
}
