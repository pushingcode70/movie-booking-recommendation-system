package models

type Seat struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	ScreenID   uint   `json:"screen_id"`
	SeatNumber string `json:"seat_number"`
	SeatType   string `json:"seat_type"`

	Screen Screen `gorm:"foreignKey:ScreenID" json:"-"`
}
