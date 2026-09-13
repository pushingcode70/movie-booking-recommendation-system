package services

import (
	"movie-booking/models"
	"movie-booking/repositories"
)

type ScreenService struct {
	repo *repositories.ScreenRepository
}

// Constructor
func NewScreenService(repo *repositories.ScreenRepository) *ScreenService {
	return &ScreenService{repo: repo}
}

// Create Screen
func (s *ScreenService) CreateScreen(screen *models.Screen) error {
	return s.repo.CreateScreen(screen)
}

// Get Screen by ID
func (s *ScreenService) GetScreenByID(id uint) (*models.Screen, error) {
	return s.repo.GetScreenByID(id)
}

// Get All Screens
func (s *ScreenService) GetAllScreens() ([]models.Screen, error) {
	return s.repo.GetAllScreens()
}

// Update Screen
func (s *ScreenService) UpdateScreen(screen *models.Screen) error {
	return s.repo.UpdateScreen(screen)
}

// Delete Screen
func (s *ScreenService) DeleteScreen(id uint) error {

	return s.repo.DeleteScreen(id)
}
