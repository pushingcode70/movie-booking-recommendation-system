package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

func (r *MovieRepository) CreateMovie(movie *models.Movie) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		// keep the genres separate while creating the movie itself
		genres := movie.Genres
		movie.Genres = nil

		// create the movie record first
		if err := tx.Create(movie).Error; err != nil {
			return err
		}

		// associate each genre with the movie
		for _, genre := range genres {

			var existingGenre models.Genre

			// reuse genre if it already exists in db
			err := tx.
				Where("tmdb_id = ?", genre.TMDBID).
				First(&existingGenre).Error

			if err != nil {

				// if genre does not exist so create
				if err := tx.Create(&genre).Error; err != nil {
					return err
				}

				existingGenre = genre
			}

			// create a movie-genre relationship
			if err := tx.Model(movie).Association("Genres").Append(&existingGenre); err != nil {
				return err
			}
		}
		return nil

	})
}

func (r *MovieRepository) GetMovieByID(id uint) (*models.Movie, error) {
	var movie models.Movie

	err := r.db.
		Preload("Genres").
		First(&movie, id).Error
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (r *MovieRepository) GetMovieByTMDBID(tmdbID int) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.Preload("Genres").Where("tmdb_id = ?", tmdbID).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie

	err := r.db.Preload("Genres").Limit(10).Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) UpdateMovie(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

func (r *MovieRepository) DeleteMovie(id uint) error {
	return r.db.Delete(&models.Movie{}, id).Error
}

func (r *MovieRepository) GetGenreByTMDBID(tmdbID int) (*models.Genre, error) {
	var genre models.Genre

	err := r.db.Where("tmdb_id = ?", tmdbID).First(&genre).Error
	if err != nil {
		return nil, err
	}

	return &genre, nil
}

func (r *MovieRepository) CreateGenre(genre *models.Genre) error {
	return r.db.Create(genre).Error
}

func (r *MovieRepository) GetMoviesByGenreID(
	genreID int,
	limit int,
) ([]models.Movie, error) {

	var movies []models.Movie

	err := r.db.
		Model(&models.Movie{}).
		Distinct("movies.*").
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Joins("JOIN genres ON genres.id = movie_genres.genre_id").
		Where("genres.tmdb_id = ?", genreID).
		Preload("Genres").
		Limit(limit).
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetMoviesByTMDBIDs(tmdbIDs []int) ([]models.Movie, error) {

	var movies []models.Movie

	if len(tmdbIDs) == 0 {
		return movies, nil
	}

	err := r.db.
		Where("tmdb_id IN ?", tmdbIDs).
		Preload("Genres").
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) GetMoviesByGenreIDs(genreIDs []int, limit int) ([]models.Movie, error) {

	var movies []models.Movie

	if len(genreIDs) == 0 {
		return movies, nil
	}

	err := r.db.
		Model(&models.Movie{}).
		Distinct("movies.*").
		Joins("JOIN movie_genres ON movie_genres.movie_id = movies.id").
		Where("movie_genres.genre_id IN ?", genreIDs).
		Preload("Genres").
		Limit(limit).
		Find(&movies).Error

	if err != nil {
		return nil, err
	}

	return movies, nil
}
