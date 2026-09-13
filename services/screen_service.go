package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
)

type ScreenService struct {
	repo *repositories.ScreenRepository
}

func NewScreenService(repo *repositories.ScreenRepository) *ScreenService {
	return &ScreenService{repo: repo}
}

func (s *ScreenService) CreateScreen(screen *models.Screen) error {
	return s.repo.CreateScreen(screen)
}

func (s *ScreenService) GetScreenByID(id uint) (*models.Screen, error) {
	return s.repo.GetScreenByID(id)
}

func (s *ScreenService) GetAllScreens() ([]models.Screen, error) {
	return s.repo.GetAllScreens()
}

func (s *ScreenService) UpdateScreen(screen *models.Screen) error {
	return s.repo.UpdateScreen(screen)
}

func (s *ScreenService) DeleteScreen(id uint) error {

	return s.repo.DeleteScreen(id)
}
