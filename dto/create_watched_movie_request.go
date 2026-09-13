package dto

type CreateWatchedMovieRequest struct {
	TMDBID int `json:"tmdb_id" binding:"required"`

	Rating *int `json:"rating"`

	Review *string `json:"review"`
}
