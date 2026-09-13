package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
)

type GenreService struct {
	repo *repositories.GenreRepository
}

func NewGenreService(repo *repositories.GenreRepository) *GenreService {
	return &GenreService{
		repo: repo,
	}
}

func (s *GenreService) GetAllGenres() ([]models.Genre, error) {
	return s.repo.GetAllGenres()
}
