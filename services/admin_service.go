package services

import (
	"movie-booking/dto"
	"movie-booking/repositories"
)

type AdminService struct {
	repo *repositories.AdminRepository
}

func NewAdminService(repo *repositories.AdminRepository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (s *AdminService) GetDashboardStats() (*dto.DashboardResponse, error) {
	return s.repo.GetDashboardStats()
}

// GetRecentPayments returns all payments, optionally filtered by a date (YYYY-MM-DD).
func (s *AdminService) GetRecentPayments(date string) ([]dto.AdminRecentPaymentResponse, error) {
	return s.repo.GetRecentPayments(date)
}

// GetRecentBookings returns bookings optionally filtered by date (YYYY-MM-DD).
// An empty string returns all bookings.
func (s *AdminService) GetRecentBookings(date string) ([]dto.AdminRecentBookingResponse, error) {
	return s.repo.GetRecentBookings(date)
}

func (s *AdminService) GetRunningShows() ([]dto.AdminRunningShowResponse, error) {
	return s.repo.GetRunningShows()
}

func (s *AdminService) GetTheatreSchedule(
	theatreID uint,
	date string,
) ([]dto.AdminScheduleResponse, error) {

	return s.repo.GetTheatreSchedule(theatreID, date)
}
