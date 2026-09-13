package dto

import "time"

type WatchedMovieResponse struct {
	TMDBID int `json:"tmdb_id"`

	Title string `json:"title"`

	Overview string `json:"overview"`

	PosterPath string `json:"poster_path"`

	ReleaseDate string `json:"release_date"`

	Rating *int `json:"rating"`

	Review *string `json:"review"`

	WatchedAt time.Time `json:"watched_at"`
}
