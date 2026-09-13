package repositories

import (
	"movie-booking/models"

	"time"

	"gorm.io/gorm"
)

type ShowRepository struct {
	db *gorm.DB
}

// Constructor
func NewShowRepository(db *gorm.DB) *ShowRepository {
	return &ShowRepository{db: db}
}

func (r *ShowRepository) WithTx(tx *gorm.DB) *ShowRepository {
	return &ShowRepository{
		db: tx,
	}
}

func (r *ShowRepository) CreateShow(show *models.Show) error {

	if err := r.db.Create(show).Error; err != nil {
		return err
	}

	return r.db.
		Preload("Movie").
		Preload("Screen").
		First(show, show.ID).Error
}

// Get Show by ID
func (r *ShowRepository) GetShowByID(id uint) (*models.Show, error) {
	var show models.Show

	err := r.db.
		Preload("Movie").
		Preload("Screen").
		First(&show, id).Error
	if err != nil {
		return nil, err
	}

	return &show, nil
}

// Get All Shows
func (r *ShowRepository) GetAllShows() ([]models.Show, error) {
	var shows []models.Show

	err := r.db.
		Preload("Movie").
		Preload("Screen").Find(&shows).Error
	if err != nil {
		return nil, err
	}

	return shows, nil
}

// Update Show
func (r *ShowRepository) UpdateShow(show *models.Show) error {
	return r.db.Save(show).Error
}

// Delete Show
func (r *ShowRepository) DeleteShow(id uint) error {
	return r.db.Delete(&models.Show{}, id).Error
}

// get activeshows
func (r *ShowRepository) GetActiveShows() ([]models.Show, error) {

	var shows []models.Show

	err := r.db.
		Preload("Movie").
		Preload("Screen").
		Where("end_time > NOW()").
		Find(&shows).Error

	return shows, err
}

func (r *ShowRepository) HasOverlappingShow(screenID uint, startTime, endTime time.Time) (bool, error) {

	var count int64

	err := r.db.Model(&models.Show{}).
		Where("screen_id = ?", screenID).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&count).Error

	return count > 0, err

}

// HasBookings checks whether a show has any bookings.
func (r *ShowRepository) HasBookings(showID uint) (bool, error) {

	var count int64

	err := r.db.
		Model(&models.Booking{}).
		Where("show_id = ?", showID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// HasFutureShows checks whether a screen has any upcoming or ongoing shows.
func (r *ShowRepository) HasFutureShows(screenID uint) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.Show{}).
		Where("screen_id = ?", screenID).
		Where("end_time > NOW()").
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
