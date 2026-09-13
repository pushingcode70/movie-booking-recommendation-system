package repositories

import (
	"movie-booking/dto"
	"movie-booking/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{
		db: db,
	}
}

func (r *AdminRepository) GetDashboardStats() (*dto.DashboardResponse, error) {
	var stats dto.DashboardResponse

	// Movies
	if err := r.db.Model(&models.Movie{}).Count(&stats.TotalMovies).Error; err != nil {
		return nil, err
	}

	// Theatres
	if err := r.db.Model(&models.Theatre{}).Count(&stats.TotalTheatres).Error; err != nil {
		return nil, err
	}

	// Screens
	if err := r.db.Model(&models.Screen{}).Count(&stats.TotalScreens).Error; err != nil {
		return nil, err
	}

	// Active Shows (end_time in the future)
	if err := r.db.Model(&models.Show{}).Where("end_time >= NOW()").Count(&stats.TotalShows).Error; err != nil {
		return nil, err
	}

	// Bookings
	if err := r.db.Model(&models.Booking{}).Count(&stats.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Successful Payments
	if err := r.db.Model(&models.Payment{}).
		Where("status = ?", "SUCCESS").
		Count(&stats.SuccessfulPayments).Error; err != nil {
		return nil, err
	}

	// Pending Payments
	if err := r.db.Model(&models.Payment{}).
		Where("status = ?", "PENDING").
		Count(&stats.PendingPayments).Error; err != nil {
		return nil, err
	}

	// Total Revenue
	if err := r.db.Model(&models.Payment{}).
		Where("status = ?", "SUCCESS").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&stats.TotalRevenue).Error; err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *AdminRepository) GetRecentPayments(date string) ([]dto.AdminRecentPaymentResponse, error) {

	var payments []dto.AdminRecentPaymentResponse

	q := r.db.
		Table("payments").
		Select(`
			payments.id AS payment_id,
			payments.booking_id AS booking_id,
			payments.amount AS amount,
			payments.status AS status,
			payments.created_at AS created_at
		`)

	if date != "" {
		q = q.Where("DATE(payments.created_at) = ?", date)
	}

	err := q.Order("payments.created_at DESC").Scan(&payments).Error
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *AdminRepository) GetRecentBookings(date string) ([]dto.AdminRecentBookingResponse, error) {

	var bookings []dto.AdminRecentBookingResponse

	q := r.db.
		Table("bookings").
		Select(`
			bookings.id AS booking_id,
			COALESCE(users.name, 'Unknown User') AS customer_name,
			COALESCE(movies.title, 'Movie unavailable') AS movie_title,
			bookings.total_amount AS amount,
			bookings.status,
			bookings.created_at,
			shows.id AS show_id,
			shows.start_time AS show_start_time,
			shows.end_time AS show_end_time
		`).
		Joins("LEFT JOIN users ON users.id = bookings.user_id").
		Joins("LEFT JOIN shows ON shows.id = bookings.show_id").
		Joins("LEFT JOIN movies ON movies.id = shows.movie_id")

	if date != "" {
		q = q.Where("DATE(bookings.created_at) = ?", date)
	}

	err := q.Order("bookings.created_at DESC").Scan(&bookings).Error
	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// GetRunningShows returns all shows that are currently running.
func (r *AdminRepository) GetRunningShows() ([]dto.AdminRunningShowResponse, error) {

	var shows []dto.AdminRunningShowResponse

	err := r.db.
		Table("shows").
		Select(`
			shows.id AS show_id,
			movies.title AS movie_title,
			theatres.name AS theatre_name,
			screens.screen_number,
			shows.start_time,
			shows.end_time
		`).
		Joins("JOIN movies ON movies.id = shows.movie_id").
		Joins("JOIN screens ON screens.id = shows.screen_id").
		Joins("JOIN theatres ON theatres.id = screens.theatre_id").
		Where("shows.start_time <= NOW()").
		Where("shows.end_time >= NOW()").
		Scan(&shows).Error

	if err != nil {
		return nil, err
	}
	// For every running show, calculate:
	// 1. Total seats in the screen
	// 2. Booked seats
	// 3. Occupancy percentage
	for i := range shows {

		// Count total seats in this screen
		var totalSeats int64

		err = r.db.
			Table("seats").
			Joins("JOIN shows ON shows.screen_id = seats.screen_id").
			Where("shows.id = ?", shows[i].ShowID).
			Count(&totalSeats).Error

		if err != nil {
			return nil, err
		}

		shows[i].TotalSeats = totalSeats

		// Count booked seats for this show
		var bookedSeats int64

		err = r.db.
			Table("booking_seats").
			Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
			Where("bookings.show_id = ?", shows[i].ShowID).
			Where("bookings.status = ?", "CONFIRMED").
			Count(&bookedSeats).Error

		if err != nil {
			return nil, err
		}

		shows[i].BookedSeats = bookedSeats

		// Calculate occupancy percentage
		if totalSeats > 0 {
			shows[i].Occupancy = float64(bookedSeats) / float64(totalSeats) * 100
		}
	}

	return shows, nil
}

// GetTheatreSchedule returns all shows for a theatre.
func (r *AdminRepository) GetTheatreSchedule(theatreID uint, date string) ([]dto.AdminScheduleResponse, error) {

	var schedule []dto.AdminScheduleResponse

	err := r.db.
		Table("shows").
		Select(`
			shows.id AS show_id,
			movies.id AS movie_id,
			movies.title AS movie_title,
			screens.id AS screen_id,
			screens.screen_number,
			shows.start_time,
			shows.end_time,
			shows.price
		`).
		Joins("JOIN movies ON movies.id = shows.movie_id").
		Joins("JOIN screens ON screens.id = shows.screen_id").
		Where("screens.theatre_id = ?", theatreID).
		Where("DATE(shows.start_time) = ?", date).
		Order("screens.screen_number ASC").
		Order("shows.start_time ASC").
		Scan(&schedule).Error

	if err != nil {
		return nil, err
	}

	return schedule, nil
}
