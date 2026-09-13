package repositories

import (
	"movie-booking/models"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieEmbeddingRepository struct {
	db *gorm.DB
}

type SimilarMovie struct {
	TMDBID   int     `gorm:"column:tmdb_id"`
	Distance float64 `gorm:"column:distance"`
}

// constructor
func NewMovieEmbeddingRepository(db *gorm.DB) *MovieEmbeddingRepository {
	return &MovieEmbeddingRepository{
		db: db,
	}
}

// Create stores a new movie embedding.
func (r *MovieEmbeddingRepository) Create(embedding *models.MovieEmbedding) error {
	return r.db.Create(embedding).Error
}

// GetByTMDBID returns the embedding for a TMDB movie.
func (r *MovieEmbeddingRepository) GetByTMDBID(tmdbID int) (*models.MovieEmbedding, error) {
	var embedding models.MovieEmbedding

	err := r.db.
		Where("tmdb_id = ?", tmdbID).
		First(&embedding).Error

	if err != nil {
		return nil, err
	}

	return &embedding, nil
}

// FindSimilarMovies returns movies whose embeddings are closest
// to the supplied query vector using cosine distance.
func (r *MovieEmbeddingRepository) FindSimilarMovies(
	queryVector pgvector.Vector,
	limit int,
	excludeTMDBIDs []int,
) ([]SimilarMovie, error) {

	var results []SimilarMovie

	query := r.db.
		Model(&models.MovieEmbedding{}).
		Select("tmdb_id, embedding <=> ? AS distance", queryVector).
		Clauses(clause.OrderBy{
			Expression: gorm.Expr("embedding <=> ?", queryVector),
		}).
		Limit(limit)

	if len(excludeTMDBIDs) > 0 {
		query = query.Where("tmdb_id NOT IN ?", excludeTMDBIDs)
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
