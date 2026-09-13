package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type BookingSeatRepository struct {
	db *gorm.DB
}

// Constructor
func NewBookingSeatRepository(db *gorm.DB) *BookingSeatRepository {
	return &BookingSeatRepository{
		db: db,
	}
}
func (r *BookingSeatRepository) WithTx(tx *gorm.DB) *BookingSeatRepository {
	return &BookingSeatRepository{
		db: tx,
	}
}

// Create BookingSeat
func (r *BookingSeatRepository) CreateBookingSeat(bookingSeat *models.BookingSeat) error {
	return r.db.Create(bookingSeat).Error
}

// Get BookingSeats By Booking ID
func (r *BookingSeatRepository) GetBookingSeatsByBookingID(bookingID uint) ([]models.BookingSeat, error) {
	var bookingSeats []models.BookingSeat

	err := r.db.Where("booking_id = ?", bookingID).Find(&bookingSeats).Error
	if err != nil {
		return nil, err
	}

	return bookingSeats, nil
}

/*input: Seat IDs the user wants to book.
Process: Check whether those IDs already exist in booking_seats.
Output: A list of the requested seat IDs that are already booked.*/
//the below func will demonstrate this

func (r *BookingSeatRepository) GetBookedSeatIDs(showID uint, seatIDs []uint) ([]uint, error) {
	var bookedSeatIDs []uint

	err := r.db.
		Model(&models.BookingSeat{}).
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("bookings.show_id = ? AND booking_seats.seat_id IN ? AND bookings.status = ?", showID, seatIDs, "CONFIRMED").
		Pluck("booking_seats.seat_id", &bookedSeatIDs).Error //only return booked seat ids not entire row

	if err != nil {
		return nil, err
	}
	return bookedSeatIDs, nil
}
