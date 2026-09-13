package repositories

import (
	"movie-booking/dto"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

// constructor
func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{
		db: db,
	}
}

// GetTicketEmailData fetches all the information
// required to send a movie ticket email.
func (r *TicketRepository) GetTicketEmailData(bookingID uint) (*dto.TicketEmailData, error) {

	var ticket dto.TicketEmailData

	// 1. Fetch Booking
	// 2. Fetch User
	// 3. Fetch Show
	// 4. Fetch Movie
	// 5. Fetch Screen
	// 6. Fetch Theatre
	// 7. Fetch Booking Seats
	// 8. Fetch Seat Numbers
	// 9. Populate TicketEmailData

	return &ticket, nil
}
