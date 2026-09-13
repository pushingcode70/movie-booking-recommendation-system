package models

type BookingSeat struct {
	/*This tells PostgreSQL:
	  For the SAME show,
	 the SAME seat can exist only once.*/
	ID uint `gorm:"primaryKey" json:"id"`

	BookingID uint `json:"booking_id"`
	SeatID    uint ` json:"seat_id"`
}
