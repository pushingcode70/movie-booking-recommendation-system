package dto

// contains the user's temporary
// inputs for generating personalized recommendations.
type RecommendationRequest struct {
	GenreIDs []int  `json:"genre_ids"`
	MovieIDs []int  `json:"movie_ids"`
	Prompt   string `json:"prompt"`
}

// represents a movie returned by the recommendation engine.
type MovieRecommendation struct {
	TMDBID   int     `json:"tmdb_id"`
	Distance float64 `json:"distance"`
}
