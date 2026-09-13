package dto

type CreateMovieRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Duration    int    `json:"duration" binding:"required"`
	Language    string `json:"language"`
	Genre       string `json:"genre"`
	Rating      string `json:"rating"`
	PosterURL   string `json:"poster_url"`
}

type UpdateMovieRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Language    string `json:"language"`
	Genre       string `json:"genre"`
	Rating      string `json:"rating"`
	PosterURL   string `json:"poster_url"`
}

type MovieResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Language    string `json:"language"`
	Genre       string `json:"genre"`
	Rating      string `json:"rating"`
	PosterURL   string `json:"poster_url"`
}
