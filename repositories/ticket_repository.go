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

// getTicketEmailData fetches all information required to send a movie ticket email
func (r *TicketRepository) GetTicketEmailData(bookingID uint) (*dto.TicketEmailData, error) {

	var ticket dto.TicketEmailData

	// 1. fetch booking
	// 2. fetch user
	// 3. fetch show
	// 4. fetch movie
	// 5. fetch screen
	// 6. fetch theatre
	// 7. fetch booking seats
	// 8. fetch seat numbers
	// 9. populate ticketEmailData

	return &ticket, nil
}
