package models

type BookingSeat struct {
	// ensures for the same show, the same seat can exist only once
	ID uint `gorm:"primaryKey" json:"id"`

	BookingID uint `json:"booking_id"`
	SeatID    uint `json:"seat_id"`
}
