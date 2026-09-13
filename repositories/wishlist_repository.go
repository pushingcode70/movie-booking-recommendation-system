package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type WishlistRepository struct {
	db *gorm.DB
}

func NewWishlistRepository(db *gorm.DB) *WishlistRepository {
	return &WishlistRepository{
		db: db,
	}
}

// create adds a movie to the user's wishlist
func (r *WishlistRepository) Create(wishlist *models.Wishlist) error {
	return r.db.Create(wishlist).Error
}

// getByUserID returns all wishlist entries for a user
func (r *WishlistRepository) GetByUserID(userID uint) ([]models.Wishlist, error) {
	var wishlist []models.Wishlist

	err := r.db.
		Where("user_id = ?", userID).
		Find(&wishlist).Error

	return wishlist, err
}

func (r *WishlistRepository) Delete(userID uint, tmdbID int) error {
	return r.db.
		Where("user_id = ? AND tmdb_id = ?", userID, tmdbID).
		Delete(&models.Wishlist{}).Error
}

func (r *WishlistRepository) Exists(userID uint, tmdbID int) (bool, error) {
	var count int64

	err := r.db.
		Model(&models.Wishlist{}).
		Where("user_id = ? AND tmdb_id = ?", userID, tmdbID).
		Count(&count).Error

	return count > 0, err
}
