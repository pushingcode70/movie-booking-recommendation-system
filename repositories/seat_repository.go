package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"

	"movie-booking/dto"
)

type SeatRepository struct {
	db *gorm.DB
}

// Constructor
func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) WithTx(tx *gorm.DB) *SeatRepository {
	return &SeatRepository{
		db: tx,
	}
}

// Create Seat
func (r *SeatRepository) CreateSeat(seat *models.Seat) error {
	return r.db.Create(seat).Error
}

// Get Seat by ID
func (r *SeatRepository) GetSeatByID(id uint) (*models.Seat, error) {
	var seat models.Seat

	err := r.db.First(&seat, id).Error
	if err != nil {
		return nil, err
	}

	return &seat, nil
}

// Get All Seats
func (r *SeatRepository) GetAllSeats() ([]models.Seat, error) {
	var seats []models.Seat

	err := r.db.Find(&seats).Error
	if err != nil {
		return nil, err
	}

	return seats, nil
}

// Update Seat
func (r *SeatRepository) UpdateSeat(seat *models.Seat) error {
	return r.db.Save(seat).Error
}

// Delete Seat
func (r *SeatRepository) DeleteSeat(id uint) error {
	return r.db.Delete(&models.Seat{}, id).Error
}

// get seats by IDs
func (r *SeatRepository) GetSeatByIDs(ids []uint) ([]models.Seat, error) {
	var seats []models.Seat

	err := r.db.Where("id IN ?", ids).Find(&seats).Error //we only need to find seats from the row of db
	/*Why IN?

	Suppose:

	ids := []uint{1, 5, 9}
	GORM generates a query similar to:
	SELECT * FROM seats
	WHERE id IN (1, 5, 9);
	Instead of:
	SELECT * FROM seats WHERE id = 1;
	SELECT * FROM seats WHERE id = 5;
	SELECT * FROM seats WHERE id = 9;*/

	if err != nil {
		return nil, err
	}

	return seats, nil
}

// GetShowSeatLayout returns all seats of the show's screen
// and marks whether each seat is already booked.
func (r *SeatRepository) GetShowSeatLayout(showID uint) ([]dto.ShowSeatRow, error) {

	var rows []dto.ShowSeatRow

	err := r.db.
		Table("seats").

		// Select seat information along with show, movie and theatre details.
		Select(`
			shows.id AS show_id,
			movies.title AS movie_title,
			theatres.name AS theatre_name,
			screens.screen_number,

			seats.id AS seat_id,
			seats.seat_number,

			CASE
				WHEN bookings.id IS NULL THEN false
				ELSE true
			END AS is_booked
		`).

		// Seat belongs to a screen.
		Joins("JOIN screens ON screens.id = seats.screen_id").

		// The requested show is running on this screen.
		Joins("JOIN shows ON shows.screen_id = screens.id").

		// Show is for this movie.
		Joins("JOIN movies ON movies.id = shows.movie_id").

		// Screen belongs to this theatre.
		Joins("JOIN theatres ON theatres.id = screens.theatre_id").

		// Find booked seats for THIS show.
		Joins(`
			LEFT JOIN booking_seats
			ON booking_seats.seat_id = seats.id
			AND booking_seats.booking_id IN (
				SELECT id FROM bookings WHERE show_id = ? AND status = 'CONFIRMED'
			)
		`, showID).
		Joins(`
			LEFT JOIN bookings
			ON bookings.id = booking_seats.booking_id
		`).

		// Requested show only.
		Where("shows.id = ?", showID).

		// Keep seats ordered naturally.

		/*A1
		A2
		A3
		...
		A10
		A11

		B1
		B2
		...instead of the alphabetical order:
		A1
		A10
		A11
		A2*/
		Order(`SUBSTRING(seats.seat_number, 1, 1),CAST(SUBSTRING(seats.seat_number FROM 2) AS INTEGER)`).Scan(&rows).Error //// Execute the query and map each returned SQL row into the ShowSeatRow slice.

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// DeleteSeatsByScreen removes all seats belonging to a screen.
func (r *SeatRepository) DeleteSeatsByScreen(screenID uint) error {
	return r.db.
		Where("screen_id = ?", screenID).
		Delete(&models.Seat{}).Error
}

// BulkCreateSeats inserts multiple seats in a single query.
func (r *SeatRepository) BulkCreateSeats(seats []models.Seat) error {
	return r.db.Create(&seats).Error
}
