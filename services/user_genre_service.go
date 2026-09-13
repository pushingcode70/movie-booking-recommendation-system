package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
)

type UserGenreService struct {
	repo *repositories.UserGenreRepository
}

func NewUserGenreService(repo *repositories.UserGenreRepository) *UserGenreService {
	return &UserGenreService{
		repo: repo,
	}
}

func (s *UserGenreService) GetFavoriteGenres(userID uint) ([]models.Genre, error) {
	return s.repo.GetByUserID(userID)
}

func (s *UserGenreService) AddFavoriteGenre(userID uint, tmdbGenreID int) error {
	return s.repo.AddGenre(userID, tmdbGenreID)
}

func (s *UserGenreService) RemoveFavoriteGenre(userID uint, tmdbGenreID int) error {
	return s.repo.RemoveGenre(userID, tmdbGenreID)
}
