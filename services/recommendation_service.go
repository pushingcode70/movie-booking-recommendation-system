package services

import (
	"errors"
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

type RecommendationService struct {
	userGenreRepo         *repositories.UserGenreRepository
	wishlistRepo          *repositories.WishlistRepository
	watchedRepo           *repositories.WatchedMovieRepository
	movieRepo             *repositories.MovieRepository
	movieEmbeddingService *MovieEmbeddingService
	movieEmbeddingRepo    *repositories.MovieEmbeddingRepository
}

func NewRecommendationService(
	userGenreRepo *repositories.UserGenreRepository, wishlistRepo *repositories.WishlistRepository,
	watchedRepo *repositories.WatchedMovieRepository, movieRepo *repositories.MovieRepository, movieEmbeddingService *MovieEmbeddingService,
	movieEmbeddingRepo *repositories.MovieEmbeddingRepository) *RecommendationService {
	return &RecommendationService{
		userGenreRepo:         userGenreRepo,
		wishlistRepo:          wishlistRepo,
		watchedRepo:           watchedRepo,
		movieRepo:             movieRepo,
		movieEmbeddingService: movieEmbeddingService,
		movieEmbeddingRepo:    movieEmbeddingRepo,
	}
}

func (s *RecommendationService) GetFavoriteGenreContext(userID uint) ([]models.Genre, error) {
	return s.userGenreRepo.GetByUserID(userID)
}

func (s *RecommendationService) GetTasteContext(userID uint) (wishlistMovieIDs []int, watchedMovieIDs []int, err error) {

	// get the user's wishlist
	wishlist, err := s.wishlistRepo.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}
	// get the user's watched list
	watchedMovies, err := s.watchedRepo.GetByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	// extract tmdb ids from both lists
	wishlistMovieIDs = make([]int, 0, len(wishlist))

	for _, movie := range wishlist {
		wishlistMovieIDs = append(wishlistMovieIDs, movie.TMDBID)
	}

	watchedMovieIDs = make([]int, 0, len(watchedMovies))

	for _, movie := range watchedMovies {
		watchedMovieIDs = append(watchedMovieIDs, movie.TMDBID)
	}

	return wishlistMovieIDs, watchedMovieIDs, nil
}

func (s *RecommendationService) GenerateAndStoreMovieEmbedding(tmdbID int) error {

	// check if it exists
	_, err := s.movieEmbeddingRepo.GetByTMDBID(tmdbID)

	if err == nil {
		return nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// if not then create it
	embedding, err := s.movieEmbeddingService.CreateMovieEmbeddingByTMDBID(tmdbID)
	if err != nil {
		return err
	}

	// create embedding
	movieEmbedding := &models.MovieEmbedding{
		TMDBID:    tmdbID,
		Embedding: pgvector.NewVector(embedding),
	}

	return s.movieEmbeddingRepo.Create(movieEmbedding)
}

func (s *RecommendationService) GetPromptRecommendations(prompt string, limit int) ([]dto.MovieRecommendation, error) {

	queryText := "Represent this sentence for searching relevant passages: " + prompt

	queryVector, err := s.movieEmbeddingService.GenerateEmbedding(queryText)

	if err != nil {
		return nil, err
	}

	results, err := s.movieEmbeddingRepo.FindSimilarMovies(pgvector.NewVector(queryVector), limit, nil)

	if err != nil {
		return nil, err
	}

	recommendations := make([]dto.MovieRecommendation, 0, len(results))

	for _, result := range results {
		recommendations = append(recommendations, dto.MovieRecommendation{
			TMDBID:   result.TMDBID,
			Distance: result.Distance,
		})
	}

	return recommendations, nil

}

func (s *RecommendationService) GetPromptGenreRecommendations(prompt string, genreIDs []int, limit int) ([]models.Movie, error) {

	if prompt == "" || len(genreIDs) == 0 || limit <= 0 {
		return []models.Movie{}, nil
	}

	// convert prompt into bge query vector embedding
	queryText := "Represent this sentence for searching relevant passages: " + prompt

	queryVector, err := s.movieEmbeddingService.GenerateEmbedding(queryText)

	if err != nil {
		return nil, err
	}

	// ask pgvector for more movies than we finally need
	// this gives us enough candidates to balance the selected genres
	candidateLimit := limit * 10

	similarMovies, err := s.movieEmbeddingRepo.FindSimilarMovies(
		pgvector.NewVector(queryVector),
		candidateLimit,
		nil,
	)

	if err != nil {
		return nil, err
	}

	// extract the tmdb ids returned by pgvector
	tmdbIDs := make([]int, 0, len(similarMovies))

	for _, result := range similarMovies {
		tmdbIDs = append(tmdbIDs, result.TMDBID)
	}

	// use tmdb ids to get the local movie objects
	movies, err := s.movieRepo.GetMoviesByTMDBIDs(tmdbIDs)

	if err != nil {
		return nil, err
	}

	// build a lookup table using tmdb id
	movieByTMDBID := make(map[int]models.Movie, len(movies))

	for _, movie := range movies {
		movieByTMDBID[movie.TMDBID] = movie
	}

	// create a separate candidate list for each selected genre
	genreMovies := make(map[int][]models.Movie)

	for _, genreID := range genreIDs {
		genreMovies[genreID] = []models.Movie{}
	}

	// put each movie into the selected genre lists it belongs to
	for _, result := range similarMovies {

		movie, exists := movieByTMDBID[result.TMDBID]

		if !exists {
			continue
		}

		// check which selected genres this movie belongs to
		for _, genre := range movie.Genres {

			if _, selected := genreMovies[genre.TMDBID]; selected {
				genreMovies[genre.TMDBID] = append(
					genreMovies[genre.TMDBID],
					movie,
				)
			}
		}
	}

	// keep the cursor for every genre
	genreIndexes := make(map[int]int)

	for _, genreID := range genreIDs {
		genreIndexes[genreID] = 0
	}

	// tracks movies already selected
	// a movie can belong to multiple genres so this prevents duplicates
	used := make(map[uint]bool)

	selected := make([]models.Movie, 0, limit)

	// round robin through selected genres
	for len(selected) < limit {

		addedThisRound := false

		for _, genreID := range genreIDs {

			movies := genreMovies[genreID]
			index := genreIndexes[genreID]

			// skip movies that were already selected through another genre
			for index < len(movies) && (used[movies[index].ID] || (movies[index].Duration > 0 && movies[index].Duration < 90)) {
				index++
			}

			// save the updated position for this genre
			genreIndexes[genreID] = index

			// if this genre has no unused candidates left move to the next genre
			if index >= len(movies) {
				continue
			}

			movie := movies[index]

			selected = append(selected, movie)
			used[movie.ID] = true

			// move this genre's cursor forward for the next round
			genreIndexes[genreID] = index + 1

			addedThisRound = true

			if len(selected) >= limit {
				break
			}
		}

		// stop if none of the selected genres have available candidates
		if !addedThisRound {
			break
		}
	}

	return selected, nil
}

func (s *RecommendationService) GetMovieRecommendations(movieIDs []int, limit int) ([]models.Movie, error) {

	if len(movieIDs) == 0 || limit <= 0 {
		return []models.Movie{}, nil
	}

	// one candidate list for every selected movie
	movieCandidates := make(map[int][]repositories.SimilarMovie)
	movieIndexes := make(map[int]int)

	for _, movieID := range movieIDs {

		// get embedding of selected movie
		embedding, err := s.movieEmbeddingRepo.GetByTMDBID(movieID)

		if err != nil {

			// create embedding if it does not exist
			if errors.Is(err, gorm.ErrRecordNotFound) {

				err = s.GenerateAndStoreMovieEmbedding(movieID)

				if err != nil {
					return nil, err
				}

				embedding, err = s.movieEmbeddingRepo.GetByTMDBID(movieID)

				if err != nil {
					return nil, err
				}

			} else {
				return nil, err
			}
		}

		// find movies similar to this selected movie
		results, err := s.movieEmbeddingRepo.FindSimilarMovies(
			embedding.Embedding,
			limit*5,
			movieIDs,
		)

		if err != nil {
			return nil, err
		}

		movieCandidates[movieID] = results
		movieIndexes[movieID] = 0
	}

	// round-robin between selected movies
	selectedTMDBIDs := make([]int, 0, limit)
	used := make(map[int]bool)

	for len(selectedTMDBIDs) < limit {

		addedThisRound := false

		for _, movieID := range movieIDs {

			candidates := movieCandidates[movieID]
			index := movieIndexes[movieID]

			for index < len(candidates) &&
				used[candidates[index].TMDBID] {
				index++
			}

			movieIndexes[movieID] = index

			if index >= len(candidates) {
				continue
			}

			tmdbID := candidates[index].TMDBID

			selectedTMDBIDs = append(selectedTMDBIDs, tmdbID)
			used[tmdbID] = true

			movieIndexes[movieID] = index + 1
			addedThisRound = true

			if len(selectedTMDBIDs) >= limit {
				break
			}
		}

		if !addedThisRound {
			break
		}
	}

	// fetch actual local movie records
	movies, err := s.movieRepo.GetMoviesByTMDBIDs(selectedTMDBIDs)

	if err != nil {
		return nil, err
	}

	// create TMDB ID -> Movie lookup
	movieByTMDBID := make(map[int]models.Movie, len(movies))

	for _, movie := range movies {
		movieByTMDBID[movie.TMDBID] = movie
	}

	// restore recommendation order
	selected := make([]models.Movie, 0, limit)

	for _, tmdbID := range selectedTMDBIDs {

		movie, exists := movieByTMDBID[tmdbID]

		if !exists {
			continue
		}

		selected = append(selected, movie)

		if len(selected) >= limit {
			break
		}
	}

	return selected, nil
}

func (s *RecommendationService) GetGenreRecommendations(genreIDs []int, limit int) ([]models.Movie, error) {

	if len(genreIDs) == 0 || limit <= 0 {
		return []models.Movie{}, nil
	}

	// candidate limit per genre loads enough candidates to build a full balanced result
	candidateLimit := limit * 2

	genreMovies := make(map[int][]models.Movie)

	// keep the cursor for every genre
	genreIndexes := make(map[int]int)

	// fetch candidates separately for each genre
	for _, genreID := range genreIDs {

		movies, err := s.movieRepo.GetMoviesByGenreID(genreID, candidateLimit)
		if err != nil {
			return nil, err
		}

		genreMovies[genreID] = movies
		genreIndexes[genreID] = 0
	}

	selected := make([]models.Movie, 0, limit)

	// tracks movies already selected
	used := make(map[uint]bool)

	// round robin through selected genres
	for len(selected) < limit {

		addedThisRound := false

		for _, genreID := range genreIDs {

			movies := genreMovies[genreID]
			index := genreIndexes[genreID]

			for index < len(movies) && (used[movies[index].ID] || (movies[index].Duration > 0 && movies[index].Duration < 90)) {
				index++
			}

			genreIndexes[genreID] = index

			if index >= len(movies) {
				continue
			}

			movie := movies[index]

			selected = append(selected, movie)
			used[movie.ID] = true

			genreIndexes[genreID] = index + 1

			addedThisRound = true

			if len(selected) >= limit {
				break
			}
		}

		// stop if no genre has available candidates
		if !addedThisRound {
			break
		}
	}

	return selected, nil
}
