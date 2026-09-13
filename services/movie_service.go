package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
	"strings"
)

type MovieService struct {
	repo        *repositories.MovieRepository
	tmdbService *TMDBService
}

// constructor
func NewMovieService(repo *repositories.MovieRepository, tmdbService *TMDBService) *MovieService {
	return &MovieService{
		repo:        repo,
		tmdbService: tmdbService,
	}
}

// create a movie
func (s *MovieService) CreateMovie(movie *models.Movie) error {
	return s.repo.CreateMovie(movie)
}

// get a movie by local primary key id
func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	return s.repo.GetMovieByID(id)
}

// get a movie by its TMDB ID, director/cast if missing
func (s *MovieService) GetMovieByTMDBID(tmdbID int) (*models.Movie, error) {
	movie, err := s.repo.GetMovieByTMDBID(tmdbID)
	if err != nil {
		return nil, err
	}

	//  if director or cast is missing, fetch from TMDB
	if movie.Director == "" && movie.Cast == "" {
		tmdbMovie, err := s.tmdbService.GetMovieDetails(tmdbID)
		if err == nil && tmdbMovie != nil && tmdbMovie.Credits != nil {
			for _, member := range tmdbMovie.Credits.Crew {
				if member.Job == "Director" {
					movie.Director = member.Name
					break
				}
			}
			var castMembers []string
			for i, member := range tmdbMovie.Credits.Cast {
				if i >= 5 {
					break
				}
				castMembers = append(castMembers, member.Name)
			}
			movie.Cast = strings.Join(castMembers, ", ")
			_ = s.repo.UpdateMovie(movie)
		}
	}

	return movie, nil
}

// get all movies
func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	return s.repo.GetAllMovies()
}

// update an existing movie
func (s *MovieService) UpdateMovie(movie *models.Movie) error {
	return s.repo.UpdateMovie(movie)
}

// delete a movie by id
func (s *MovieService) DeleteMovie(id uint) error {
	return s.repo.DeleteMovie(id)
}
