package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type WatchedMovieRepository struct {
	db *gorm.DB
}

func NewWatchedMovieRepository(db *gorm.DB) *WatchedMovieRepository {
	return &WatchedMovieRepository{
		db: db,
	}
}

func (r *WatchedMovieRepository) Create(watched *models.WatchedMovie) error {
	return r.db.Create(watched).Error
}

func (r *WatchedMovieRepository) GetByUserID(userID uint) ([]models.WatchedMovie, error) {
	var watched []models.WatchedMovie
	err := r.db.Where("user_id = ?", userID).Find(&watched).Error
	return watched, err
}

func (r *WatchedMovieRepository) GetByUserAndTMDBID(userID uint, tmdbID int) (*models.WatchedMovie, error) {
	var watched models.WatchedMovie
	err := r.db.Where("user_id = ? AND tmdb_id = ?", userID, tmdbID).First(&watched).Error
	if err != nil {
		return nil, err
	}
	return &watched, nil
}

func (r *WatchedMovieRepository) Update(watched *models.WatchedMovie) error {
	return r.db.Save(watched).Error
}

func (r *WatchedMovieRepository) Delete(userID uint, tmdbID int) error {
	return r.db.Where("user_id = ? AND tmdb_id = ?", userID, tmdbID).Delete(&models.WatchedMovie{}).Error
}

func (r *WatchedMovieRepository) Exists(userID uint, tmdbID int) (bool, error) {
	var count int64
	err := r.db.Model(&models.WatchedMovie{}).Where("user_id = ? AND tmdb_id = ?", userID, tmdbID).Count(&count).Error
	return count > 0, err
}
