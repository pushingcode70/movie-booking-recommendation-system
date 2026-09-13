package repositories

import (
	"fmt"
	"movie-booking/models"

	"gorm.io/gorm"
)

type ScreenRepository struct {
	db *gorm.DB
}

// Constructor
func NewScreenRepository(db *gorm.DB) *ScreenRepository {
	return &ScreenRepository{db: db}
}

// Create Screen
func (r *ScreenRepository) CreateScreen(screen *models.Screen) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create the screen record
		if err := tx.Create(screen).Error; err != nil {
			return err
		}

		// 2. Generate exactly TotalSeats physical seats.
		totalSeats := screen.TotalSeats
		if totalSeats <= 0 {
			return nil
		}

		const seatsPerRow = 10

		var seats []models.Seat
		for seatIndex := 0; seatIndex < totalSeats; seatIndex++ {
			row := seatIndex / seatsPerRow
			seatNumber := (seatIndex % seatsPerRow) + 1

			rowLetter := string(rune('A' + row))

			seats = append(seats, models.Seat{
				ScreenID:   screen.ID,
				SeatNumber: fmt.Sprintf("%s%d", rowLetter, seatNumber),
				SeatType:   "REGULAR",
			})
		}

		if err := tx.Create(&seats).Error; err != nil {
			return err
		}

		return nil
	})
}

// Get Screen by ID
func (r *ScreenRepository) GetScreenByID(id uint) (*models.Screen, error) {
	var screen models.Screen

	err := r.db.First(&screen, id).Error
	if err != nil {
		return nil, err
	}

	return &screen, nil
}

// Get All Screens
func (r *ScreenRepository) GetAllScreens() ([]models.Screen, error) {
	var screens []models.Screen

	err := r.db.Find(&screens).Error
	if err != nil {
		return nil, err
	}

	return screens, nil
}

// Update Screen
func (r *ScreenRepository) UpdateScreen(screen *models.Screen) error {
	return r.db.Save(screen).Error
}

func (r *ScreenRepository) DeleteScreen(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("screen_id = ?", id).Delete(&models.Seat{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.Screen{}, id).Error; err != nil {
			return err
		}

		return nil
	})
}

// Use the repository with an existing transaction.
func (r *ScreenRepository) WithTx(tx *gorm.DB) *ScreenRepository {
	return &ScreenRepository{
		db: tx,
	}
}

// updates the total seat capacity of a screen.
func (r *ScreenRepository) UpdateTotalSeats(screenID uint, totalSeats int) error {
	return r.db.
		Model(&models.Screen{}).
		Where("id = ?", screenID).
		Update("total_seats", totalSeats).Error
}
