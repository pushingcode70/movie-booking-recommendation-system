package models

type Screen struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	TheatreID    uint `json:"theatre_id"`
	ScreenNumber int  `json:"screen_number"`
	TotalSeats   int  `json:"total_seats"`

	Theatre Theatre `gorm:"foreignKey:TheatreID" json:"-"`
}
