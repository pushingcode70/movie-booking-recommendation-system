package models

// genre represents a movie genre provided by tmdb
type Genre struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	TMDBID int    `gorm:"uniqueIndex;not null" json:"tmdb_id"`
	Name   string `gorm:"not null" json:"name"`
}

// movie represents a movie stored in the local database
type Movie struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	TMDBID int  `gorm:"uniqueIndex" json:"tmdb_id"` // tmdb movie id

	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Language    string `json:"language"`
	ReleaseDate string `json:"release_date"`

	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`

	Director string `json:"director"`
	Cast     string `json:"cast"`

	// a movie can have multiple genres and a genre can belong to multiple movies
	Genres []Genre `gorm:"many2many:movie_genres;" json:"genres"`
}
