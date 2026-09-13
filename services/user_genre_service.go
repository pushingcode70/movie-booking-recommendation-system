package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
)

type UserGenreService struct {
	repo *repositories.UserGenreRepository
}

// constructor
func NewUserGenreService(repo *repositories.UserGenreRepository) *UserGenreService {
	return &UserGenreService{
		repo: repo,
	}
}

// get all favorite genres selected by a user
func (s *UserGenreService) GetFavoriteGenres(userID uint) ([]models.Genre, error) {
	return s.repo.GetByUserID(userID)
}

// add a genre to the user's favorite genres
func (s *UserGenreService) AddFavoriteGenre(userID uint, tmdbGenreID int) error {
	return s.repo.AddGenre(userID, tmdbGenreID)
}

// remove a genre from the user's favorite genres
func (s *UserGenreService) RemoveFavoriteGenre(userID uint, tmdbGenreID int) error {
	return s.repo.RemoveGenre(userID, tmdbGenreID)
}
