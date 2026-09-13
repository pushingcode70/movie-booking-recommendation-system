package dto

type CreateWishlistRequest struct {
	TMDBID int `json:"tmdb_id" binding:"required"`
}
