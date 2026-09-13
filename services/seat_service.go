package services

import (
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
	"strings"

	"errors"
	"fmt"

	"gorm.io/gorm"
)

type SeatService struct {
	seatRepo   *repositories.SeatRepository
	screenRepo *repositories.ScreenRepository
	showRepo   *repositories.ShowRepository
	db         *gorm.DB
}

// Constructor
func NewSeatService(
	db *gorm.DB,
	seatRepo *repositories.SeatRepository,
	screenRepo *repositories.ScreenRepository,
	showRepo *repositories.ShowRepository,
) *SeatService {
	return &SeatService{
		db:         db,
		seatRepo:   seatRepo,
		screenRepo: screenRepo,
		showRepo:   showRepo,
	}
}

// Create Seat
func (s *SeatService) CreateSeat(seat *models.Seat) error {
	return s.seatRepo.CreateSeat(seat)
}

// Get Seat by ID
func (s *SeatService) GetSeatByID(id uint) (*models.Seat, error) {
	return s.seatRepo.GetSeatByID(id)
}

// Get All Seats
func (s *SeatService) GetAllSeats() ([]models.Seat, error) {
	return s.seatRepo.GetAllSeats()
}

// Update Seat
func (s *SeatService) UpdateSeat(seat *models.Seat) error {
	return s.seatRepo.UpdateSeat(seat)
}

// Delete Seat
func (s *SeatService) DeleteSeat(id uint) error {
	return s.seatRepo.DeleteSeat(id)
}

// GetShowSeatLayout returns the seat layout for a show.
func (s *SeatService) GetShowSeatLayout(showID uint) (*dto.ShowSeatLayoutResponse, error) {

	// Get flat seat data from the repository.
	rows, err := s.seatRepo.GetShowSeatLayout(showID)
	if err != nil {
		return nil, err
	}

	// No seats found.
	if len(rows) == 0 {
		return nil, nil
	}

	response := &dto.ShowSeatLayoutResponse{
		ShowID:       rows[0].ShowID,
		MovieTitle:   rows[0].MovieTitle,
		TheatreName:  rows[0].TheatreName,
		ScreenNumber: rows[0].ScreenNumber,
	}

	// Map row label ("A", "B", etc.) to its position in response.Rows.
	rowIndex := make(map[string]int)

	for _, seat := range rows {

		// Extract row label from seat number.
		// Example: "A12" -> "A"
		row := strings.ToUpper(string(seat.SeatNumber[0]))

		index, exists := rowIndex[row]

		// Create a new row group if it doesn't exist.
		if !exists {
			response.Rows = append(response.Rows, dto.SeatRowGroup{
				Row:   row,
				Seats: []dto.SeatInfo{},
			})

			index = len(response.Rows) - 1
			rowIndex[row] = index
		}

		status := "available"
		if seat.IsBooked {
			status = "booked"
		}

		response.Rows[index].Seats = append(response.Rows[index].Seats, dto.SeatInfo{
			SeatID:     seat.SeatID,
			SeatNumber: seat.SeatNumber,
			Status:     status,
		})
	}

	return response, nil
}

func (s *SeatService) GenerateSeats(
	screenID uint,
	req dto.GenerateSeatsRequest,
) error {

	_, err := s.screenRepo.GetScreenByID(screenID)
	if err != nil {
		return err
	}

	hasShows, err := s.showRepo.HasFutureShows(screenID)
	if err != nil {
		return err
	}

	if hasShows {
		return errors.New("cannot regenerate seat layout while future shows exist")
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	seatRepo := s.seatRepo.WithTx(tx)
	screenRepo := s.screenRepo.WithTx(tx)

	if err := seatRepo.DeleteSeatsByScreen(screenID); err != nil {
		tx.Rollback()
		return err
	}

	var seats []models.Seat

	for row := 0; row < req.Rows; row++ {

		// Convert 0 -> A, 1 -> B, ...
		rowLetter := string(rune('A' + row))

		for seat := 1; seat <= req.SeatsPerRow; seat++ {

			seats = append(seats, models.Seat{
				ScreenID:   screenID,
				SeatNumber: fmt.Sprintf("%s%d", rowLetter, seat),
			})
		}
	}

	if err := seatRepo.BulkCreateSeats(seats); err != nil {
		tx.Rollback()
		return err
	}

	if err := screenRepo.UpdateTotalSeats(screenID, len(seats)); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
