package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) WithTx(tx *gorm.DB) *BookingRepository {
	return &BookingRepository{
		db: tx,
	}
}

func (r *BookingRepository) CreateBooking(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) GetBookingByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	err := r.db.First(&booking, id).Error
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (r *BookingRepository) GetBookingsByUserID(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking

	err := r.db.Where("user_id = ?", userID).Find(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", status).Error
}
