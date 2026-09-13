package repositories

import (
	"errors"
	"movie-booking/models"

	"gorm.io/gorm"
)

type UserGenreRepository struct {
	db *gorm.DB
}

// constructor
func NewUserGenreRepository(db *gorm.DB) *UserGenreRepository {
	return &UserGenreRepository{
		db: db,
	}
}

// getByUserID returns all favorite genres selected by the user
func (r *UserGenreRepository) GetByUserID(userID uint) ([]models.Genre, error) {
	var genres []models.Genre

	err := r.db.
		Model(&models.User{ID: userID}).
		Association("FavoriteGenres").
		Find(&genres)

	return genres, err
}

// addGenre adds a genre to the user's favorite genres
func (r *UserGenreRepository) AddGenre(userID uint, tmdbGenreID int) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var genre models.Genre
	if err := r.db.Where("tmdb_id = ?", tmdbGenreID).First(&genre).Error; err != nil {
		return err
	}

	// check if the user already has this genre as a favorite
	var count int64

	err := r.db.
		Table("user_favorite_genres").
		Where("user_id = ? AND genre_id = ?", userID, genre.ID).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("genre is already in favorite genres")
	}

	return r.db.Model(&user).
		Association("FavoriteGenres").
		Append(&genre)
}

// removeGenre removes a genre from the user's favorite genres
func (r *UserGenreRepository) RemoveGenre(userID uint, tmdbGenreID int) error {
	var user models.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return err
	}

	var genre models.Genre
	if err := r.db.Where("tmdb_id = ?", tmdbGenreID).First(&genre).Error; err != nil {
		return err
	}

	return r.db.Model(&user).
		Association("FavoriteGenres").
		Delete(&genre)
}
