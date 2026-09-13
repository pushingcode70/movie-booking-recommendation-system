package services

import (
	"errors"
	"time"

	"fmt"
	"movie-booking/models"
	"movie-booking/repositories"
)

type ShowService struct {
	showRepo   *repositories.ShowRepository
	movieRepo  *repositories.MovieRepository
	screenRepo *repositories.ScreenRepository
}

// Constructor
func NewShowService(showRepo *repositories.ShowRepository, movieRepo *repositories.MovieRepository, screenRepo *repositories.ScreenRepository) *ShowService {
	return &ShowService{
		showRepo:   showRepo,
		movieRepo:  movieRepo,
		screenRepo: screenRepo,
	}
}

// Create Show
func (s *ShowService) CreateShow(show *models.Show) error {

	//check movie exists
	movie, err := s.movieRepo.GetMovieByID(show.MovieID)
	if err != nil {
		return err
	}

	//check screen exists
	_, err = s.screenRepo.GetScreenByID(show.ScreenID)
	if err != nil {
		return err
	}

	//calculate and time
	show.EndTime = show.StartTime.Add(time.Duration(movie.Duration) * time.Minute)

	//check overlapping shows
	exists, err := s.showRepo.HasOverlappingShow(show.ScreenID, show.StartTime, show.EndTime)

	if err != nil {
		return err
	}

	if exists {
		return errors.New("show overlaps with an existing show")
	}

	return s.showRepo.CreateShow(show)
}

// Get Show by ID
func (s *ShowService) GetShowByID(id uint) (*models.Show, error) {
	return s.showRepo.GetShowByID(id)
}

// Get All Shows
func (s *ShowService) GetAllShows() ([]models.Show, error) {
	return s.showRepo.GetActiveShows()
}

// Update Show
func (s *ShowService) UpdateShow(show *models.Show) error {

	movie, err := s.movieRepo.GetMovieByID(show.MovieID)

	if err != nil {
		return err
	}

	show.EndTime = show.StartTime.Add(
		time.Duration(movie.Duration) * time.Minute,
	)

	fmt.Println("Movie Duration:", movie.Duration)
	fmt.Println("Start Time:", show.StartTime)
	fmt.Println("End Time:", show.EndTime)

	return s.showRepo.UpdateShow(show)
}

func (s *ShowService) DeleteShow(id uint) error {

	// Check if the show has any bookings.
	hasBookings, err := s.showRepo.HasBookings(id)
	if err != nil {
		return err
	}

	if hasBookings {
		return errors.New("cannot delete show because bookings already exist")
	}

	return s.showRepo.DeleteShow(id)
}
