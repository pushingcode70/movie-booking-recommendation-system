package repositories

import (
	"movie-booking/models"

	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

// contructor
func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		db: db,
	}
}

// create a new movie
// creates a movie and assosiates its gneres
// genres are shared accross movies so an existing genre is reused instead of attempting to insert a duplicate genres
func (r *MovieRepository) CreateMovie(movie *models.Movie) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		//keep the genres separate while creating the movie itself
		genres := movie.Genres
		movie.Genres = nil

		//create the movie record first
		if err := tx.Create(movie).Error; err != nil {
			return err
		}

		//assosiate each genre with the movie.
		for _, genre := range genres {

			var existingGenre models.Genre

			//reuse genre if it already exists in db
			err := tx.
				Where("tmdb_id = ?", genre.TMDBID).
				First(&existingGenre).Error

			if err != nil {

				//if genre does not exist so create
				if err := tx.Create(&genre).Error; err != nil {
					return err
				}

				existingGenre = genre
			}

			//create a movie-genre relationship
			if err := tx.Model(movie).Association("Genres").Append(&existingGenre); err != nil {
				return err
			}
		}
		return nil

	})
}

// get movie by its local primary key id
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

// GetMovieByTMDBID finds a movie by its TMDb ID.
func (r *MovieRepository) GetMovieByTMDBID(tmdbID int) (*models.Movie, error) {
	var movie models.Movie
	err := r.db.Preload("Genres").Where("tmdb_id = ?", tmdbID).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

//get all movies

func (r *MovieRepository) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie

	err := r.db.Preload("Genres").Limit(10).Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// update an existing movie
func (r *MovieRepository) UpdateMovie(movie *models.Movie) error {
	return r.db.Save(movie).Error
}

// delete movie by id
func (r *MovieRepository) DeleteMovie(id uint) error {
	return r.db.Delete(&models.Movie{}, id).Error
}

// get genre by tmdb id
func (r *MovieRepository) GetGenreByTMDBID(tmdbID int) (*models.Genre, error) {
	var genre models.Genre

	err := r.db.Where("tmdb_id = ?", tmdbID).First(&genre).Error
	if err != nil {
		return nil, err
	}

	return &genre, nil
}

// create a new genre
func (r *MovieRepository) CreateGenre(genre *models.Genre) error {
	return r.db.Create(genre).Error
}

// GetMoviesByGenreID returns candidate movies belonging to a TMDB genre.
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

// returns candidate movies that belong to at least
// one of the supplied genres the service layer is responsible for
// balancing movies that match multiple genres and individual genres
func (r *MovieRepository) GetMoviesByGenreIDs(genreIDs []int, limit int) ([]models.Movie, error) {

	var movies []models.Movie

	if len(genreIDs) == 0 { //if user search nothing we return empty result
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
