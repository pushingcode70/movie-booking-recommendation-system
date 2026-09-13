package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type GenreRepository struct {
	db *gorm.DB
}

// constructor
func NewGenreRepository(db *gorm.DB) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

// get all genres
func (r *GenreRepository) GetAllGenres() ([]models.Genre, error) {
	var genres []models.Genre

	err := r.db.Order("name ASC").Find(&genres).Error
	if err != nil {
		return nil, err
	}

	return genres, nil
}
