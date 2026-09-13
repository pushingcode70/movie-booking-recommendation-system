package repositories

import (
	"movie-booking/dto"
	"movie-booking/models"

	"gorm.io/gorm"
)

type TheatreRepository struct {
	db *gorm.DB
}

func NewTheatreRepository(db *gorm.DB) *TheatreRepository {
	return &TheatreRepository{
		db: db,
	}
}

func (r *TheatreRepository) CreateTheatre(theatre *models.Theatre) error {
	return r.db.Create(theatre).Error
}

func (r *TheatreRepository) GetTheatreByID(id uint) (*models.Theatre, error) {
	var theatre models.Theatre

	err := r.db.First(&theatre, id).Error
	if err != nil {
		return nil, err
	}
	return &theatre, nil
}

func (r *TheatreRepository) GetAllTheatres() ([]dto.TheatreListResponse, error) {

	var theatres []dto.TheatreListResponse

	err := r.db.
		Model(&models.Theatre{}).
		Select("id, name, address").
		Order("name ASC").
		Scan(&theatres).Error

	if err != nil {
		return nil, err
	}

	return theatres, nil
}

func (r *TheatreRepository) UpdateTheatre(theatre *models.Theatre) error {
	return r.db.Save(theatre).Error
}

func (r *TheatreRepository) DeleteTheatre(id uint) error {
	return r.db.Delete(&models.Theatre{}, id).Error // use theatre model to identify table and delete record by id
}

func (r *TheatreRepository) GetTheatreSchedule(theatreID uint, date string) ([]dto.TheatreScheduleRow, error) {

	var rows []dto.TheatreScheduleRow

	err := r.db.
		Table("shows").
		Select(`
			theatres.id AS theatre_id,
			theatres.name AS theatre_name,
			theatres.address,

			movies.id AS movie_id,
			movies.title,
			movies.poster_path,
			movies.language,
			movies.duration,

			shows.id AS show_id,
			shows.start_time,
			shows.price
		`).
		Joins("JOIN screens ON screens.id = shows.screen_id").
		Joins("JOIN theatres ON theatres.id = screens.theatre_id").
		Joins("JOIN movies ON movies.id = shows.movie_id").
		Where("theatres.id = ?", theatreID).
		Where("DATE(shows.start_time) = ?", date).
		Where("shows.end_time > NOW()").
		Order("movies.title ASC").
		Order("shows.start_time ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}
