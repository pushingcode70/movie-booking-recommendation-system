package services

import (
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
)

type TheatreService struct {
	repo *repositories.TheatreRepository
}

func NewTheatreService(repo *repositories.TheatreRepository) *TheatreService {
	return &TheatreService{repo: repo}
}

func (s *TheatreService) CreateTheatre(theatre *models.Theatre) error {
	return s.repo.CreateTheatre(theatre)
}

func (s *TheatreService) GetTheatreByID(id uint) (*models.Theatre, error) {
	return s.repo.GetTheatreByID(id)
}

func (s *TheatreService) GetAllTheatres() ([]dto.TheatreListResponse, error) {
	return s.repo.GetAllTheatres()
}

func (s *TheatreService) UpdateTheatre(theatre *models.Theatre) error {
	return s.repo.UpdateTheatre(theatre)
}

func (s *TheatreService) DeleteTheatre(id uint) error {
	return s.repo.DeleteTheatre(id)
}

func (s *TheatreService) GetTheatreSchedule(
	theatreID uint,
	date string,
) (*dto.TheatreScheduleResponse, error) {

	// fetch flat rows from repository
	rows, err := s.repo.GetTheatreSchedule(theatreID, date)
	if err != nil {
		return nil, err
	}

	// if no shows found, return empty schedule instead of nil
	if len(rows) == 0 {
		return &dto.TheatreScheduleResponse{
			TheatreID: theatreID,
			Date:      date,
			Movies:    []dto.TheatreMovieSchedule{},
		}, nil
	}

	// take theatre info from first row
	response := &dto.TheatreScheduleResponse{
		TheatreID:   rows[0].TheatreID,
		TheatreName: rows[0].TheatreName,
		Address:     rows[0].Address,
		Date:        date,
		Movies:      []dto.TheatreMovieSchedule{},
	}

	// map movie id -> index in response.Movies
	movieMap := make(map[uint]int)

	// process every show returned from database
	for _, row := range rows {

		// check whether movie has been added
		index, exists := movieMap[row.MovieID]

		// first time seeing this movie
		if !exists {

			// create new movie entry
			response.Movies = append(response.Movies, dto.TheatreMovieSchedule{
				MovieID:    row.MovieID,
				Title:      row.Title,
				PosterPath: row.PosterPath,
				Language:   row.Language,
				Duration:   row.Duration,
				ShowTimes:  []dto.TheatreShowTime{},
			})

			// store index of newly added movie
			index = len(response.Movies) - 1
			movieMap[row.MovieID] = index
		}

		// add current show's timing under its movie
		response.Movies[index].ShowTimes = append(
			response.Movies[index].ShowTimes,
			dto.TheatreShowTime{
				ShowID: row.ShowID,

				// format time into user-friendly string
				StartTime: row.StartTime.Format("03:04 PM"),

				Price: row.Price,
			},
		)
	}

	// return grouped schedule
	return response, nil
}
