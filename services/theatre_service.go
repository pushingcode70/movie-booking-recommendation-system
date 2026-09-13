package services

import (
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
)

type TheatreService struct {
	repo *repositories.TheatreRepository
}

// Constructor
func NewTheatreService(repo *repositories.TheatreRepository) *TheatreService {
	return &TheatreService{repo: repo}
}

// Create Theatre
func (s *TheatreService) CreateTheatre(theatre *models.Theatre) error {
	return s.repo.CreateTheatre(theatre)
}

// Get Theatre by ID
func (s *TheatreService) GetTheatreByID(id uint) (*models.Theatre, error) {
	return s.repo.GetTheatreByID(id)
}

func (s *TheatreService) GetAllTheatres() ([]dto.TheatreListResponse, error) {
	return s.repo.GetAllTheatres()
}

// Update Theatre
func (s *TheatreService) UpdateTheatre(theatre *models.Theatre) error {
	return s.repo.UpdateTheatre(theatre)
}

// Delete Theatre
func (s *TheatreService) DeleteTheatre(id uint) error {
	return s.repo.DeleteTheatre(id)
}

// GetTheatreSchedule returns the customer-facing schedule for a theatre on a given date.
func (s *TheatreService) GetTheatreSchedule(
	theatreID uint,
	date string,
) (*dto.TheatreScheduleResponse, error) {

	// Fetch flat rows from the repository.
	// Each row represents ONE show.
	rows, err := s.repo.GetTheatreSchedule(theatreID, date)
	if err != nil {
		return nil, err
	}

	// If no shows are found, return an empty schedule instead of nil.
	// This makes frontend handling easier.
	if len(rows) == 0 {
		return &dto.TheatreScheduleResponse{
			TheatreID: theatreID,
			Date:      date,
			Movies:    []dto.TheatreMovieSchedule{},
		}, nil
	}

	// Since every row belongs to the same theatre,
	// we can take theatre information from the first row.
	response := &dto.TheatreScheduleResponse{
		TheatreID:   rows[0].TheatreID,
		TheatreName: rows[0].TheatreName,
		Address:     rows[0].Address,
		Date:        date,
		Movies:      []dto.TheatreMovieSchedule{},
	}

	// Maps MovieID -> index in response.Movies.
	// This helps us avoid looping through Movies every time
	// we want to append another showtime.
	movieMap := make(map[uint]int)

	// Process every show returned from the database.
	for _, row := range rows {

		// Check whether we've already added this movie
		// to the response.
		index, exists := movieMap[row.MovieID]

		// First time seeing this movie.
		if !exists {

			// Create a new movie entry.
			response.Movies = append(response.Movies, dto.TheatreMovieSchedule{
				MovieID:    row.MovieID,
				Title:      row.Title,
				PosterPath: row.PosterPath,
				Language:   row.Language,
				Duration:   row.Duration,
				ShowTimes:  []dto.TheatreShowTime{},
			})

			// Store the index of the newly added movie.
			index = len(response.Movies) - 1
			movieMap[row.MovieID] = index
		}

		// Add the current show's timing under its movie.
		response.Movies[index].ShowTimes = append(
			response.Movies[index].ShowTimes,
			dto.TheatreShowTime{
				ShowID: row.ShowID,

				// Convert time.Time into a user-friendly format.
				// Example:
				// 2026-07-25T19:30:00+05:30
				// becomes
				// 07:30 PM
				StartTime: row.StartTime.Format("03:04 PM"),

				Price: row.Price,
			},
		)
	}

	// Return the grouped schedule.
	return response, nil
}
