package dto

type UpdateWatchedMovieRequest struct {
	Rating *int `json:"rating"`

	Review *string `json:"review"`
}
