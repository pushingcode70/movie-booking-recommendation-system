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

func NewMovieService(repo *repositories.MovieRepository, tmdbService *TMDBService) *MovieService {
	return &MovieService{
		repo:        repo,
		tmdbService: tmdbService,
	}
}

func (s *MovieService) CreateMovie(movie *models.Movie) error {
	return s.repo.CreateMovie(movie)
}

func (s *MovieService) GetMovieByID(id uint) (*models.Movie, error) {
	return s.repo.GetMovieByID(id)
}

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

func (s *MovieService) GetAllMovies() ([]models.Movie, error) {
	return s.repo.GetAllMovies()
}

func (s *MovieService) UpdateMovie(movie *models.Movie) error {
	return s.repo.UpdateMovie(movie)
}

func (s *MovieService) DeleteMovie(id uint) error {
	return s.repo.DeleteMovie(id)
}
