package models

// Genre represents a movie genre provided by TMDB.
//
// TMDBID is kept so that the local genre can be matched with
// TMDB's genre identifier when movie metadata is synchronized.
type Genre struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	TMDBID int    `gorm:"uniqueIndex;not null" json:"tmdb_id"`
	Name   string `gorm:"not null" json:"name"`
}

// Movie represents a movie stored in the local database.
type Movie struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	TMDBID int  `gorm:"uniqueIndex" json:"tmdb_id"` // TMDb movie ID

	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Language    string `json:"language"`
	ReleaseDate string `json:"release_date"`

	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`

	Director string `json:"director"`
	Cast     string `json:"cast"`

	// A movie can have multiple genres, and a genre can belong
	// to multiple movies. GORM manages this relationship through
	// the movie_genres join table.
	Genres []Genre `gorm:"many2many:movie_genres;" json:"genres"`
}
