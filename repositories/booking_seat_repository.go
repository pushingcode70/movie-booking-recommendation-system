package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type BookingSeatRepository struct {
	db *gorm.DB
}

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

func (r *BookingSeatRepository) CreateBookingSeat(bookingSeat *models.BookingSeat) error {
	return r.db.Create(bookingSeat).Error
}

func (r *BookingSeatRepository) GetBookingSeatsByBookingID(bookingID uint) ([]models.BookingSeat, error) {
	var bookingSeats []models.BookingSeat

	err := r.db.Where("booking_id = ?", bookingID).Find(&bookingSeats).Error
	if err != nil {
		return nil, err
	}

	return bookingSeats, nil
}

func (r *BookingSeatRepository) GetBookedSeatIDs(showID uint, seatIDs []uint) ([]uint, error) {
	var bookedSeatIDs []uint

	err := r.db.
		Model(&models.BookingSeat{}).
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("bookings.show_id = ? AND booking_seats.seat_id IN ? AND bookings.status = ?", showID, seatIDs, "CONFIRMED").
		Pluck("booking_seats.seat_id", &bookedSeatIDs).Error // return only booked seat ids

	if err != nil {
		return nil, err
	}
	return bookedSeatIDs, nil
}
