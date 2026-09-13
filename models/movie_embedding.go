package models

import "github.com/pgvector/pgvector-go"

// stores the vector representation of a TMDB movie.
type MovieEmbedding struct {
	ID uint `gorm:"primaryKey"`

	//identifies the movie in TMDB.
	TMDBID int `gorm:"uniqueIndex;not null"`

	//stores the BGE embedding vector.
	//BGE base produces 768-dimensional vectors.
	Embedding pgvector.Vector `gorm:"type:vector(768);not null"`
}
