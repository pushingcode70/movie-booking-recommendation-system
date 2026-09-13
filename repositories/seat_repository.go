package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"

	"movie-booking/dto"
)

type SeatRepository struct {
	db *gorm.DB
}

func NewSeatRepository(db *gorm.DB) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) WithTx(tx *gorm.DB) *SeatRepository {
	return &SeatRepository{
		db: tx,
	}
}

func (r *SeatRepository) CreateSeat(seat *models.Seat) error {
	return r.db.Create(seat).Error
}

func (r *SeatRepository) GetSeatByID(id uint) (*models.Seat, error) {
	var seat models.Seat

	err := r.db.First(&seat, id).Error
	if err != nil {
		return nil, err
	}

	return &seat, nil
}

func (r *SeatRepository) GetAllSeats() ([]models.Seat, error) {
	var seats []models.Seat

	err := r.db.Find(&seats).Error
	if err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *SeatRepository) UpdateSeat(seat *models.Seat) error {
	return r.db.Save(seat).Error
}

func (r *SeatRepository) DeleteSeat(id uint) error {
	return r.db.Delete(&models.Seat{}, id).Error
}

func (r *SeatRepository) GetSeatByIDs(ids []uint) ([]models.Seat, error) {
	var seats []models.Seat

	err := r.db.Where("id IN ?", ids).Find(&seats).Error

	if err != nil {
		return nil, err
	}

	return seats, nil
}

func (r *SeatRepository) GetShowSeatLayout(showID uint) ([]dto.ShowSeatRow, error) {

	var rows []dto.ShowSeatRow

	err := r.db.
		Table("seats").

		// select seat information along with show, movie and theatre details
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

		// seat belongs to a screen
		Joins("JOIN screens ON screens.id = seats.screen_id").

		// requested show is running on this screen
		Joins("JOIN shows ON shows.screen_id = screens.id").

		// show is for this movie
		Joins("JOIN movies ON movies.id = shows.movie_id").

		// screen belongs to this theatre
		Joins("JOIN theatres ON theatres.id = screens.theatre_id").

		// find booked seats for this show
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

		// requested show only
		Where("shows.id = ?", showID).

		// keep seats ordered naturally
		Order(`SUBSTRING(seats.seat_number, 1, 1),CAST(SUBSTRING(seats.seat_number FROM 2) AS INTEGER)`).Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *SeatRepository) DeleteSeatsByScreen(screenID uint) error {
	return r.db.
		Where("screen_id = ?", screenID).
		Delete(&models.Seat{}).Error
}

func (r *SeatRepository) BulkCreateSeats(seats []models.Seat) error {
	return r.db.Create(&seats).Error
}
