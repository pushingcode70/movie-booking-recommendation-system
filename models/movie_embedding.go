package models

import "github.com/pgvector/pgvector-go"

// stores the vector representation of a tmdb movie
type MovieEmbedding struct {
	ID uint `gorm:"primaryKey"`

	// identifies the movie in tmdb
	TMDBID int `gorm:"uniqueIndex;not null"`

	// stores the bge embedding vector (768 dimensions)
	Embedding pgvector.Vector `gorm:"type:vector(768);not null"`
}
