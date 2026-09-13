package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"math/rand"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"movie-booking/config"
	"movie-booking/models"
	"movie-booking/repositories"

	"gorm.io/gorm"
)

// tmdb service handles communication with tmdb api
type TMDBService struct {
	apiKey    string
	client    *http.Client
	movieRepo *repositories.MovieRepository
}

// constructor
func NewTMDBService(movieRepo *repositories.MovieRepository) *TMDBService {
	return &TMDBService{
		apiKey: config.AppConfig.TMDBAPIKey,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 20,
				IdleConnTimeout:     45 * time.Second,
				TLSHandshakeTimeout: 5 * time.Second,
				DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
					dialer := &net.Dialer{
						Timeout:   5 * time.Second,
						KeepAlive: 30 * time.Second,
					}

					// Force TMDB connections to use IPv4.
					return dialer.DialContext(ctx, "tcp4", address)
				},
			},
		},
		movieRepo: movieRepo,
	}
}

// tmdb response for a movie search
type TMDBSearchResponse struct {
	Page         int         `json:"page"`
	Results      []TMDBMovie `json:"results"`
	TotalPages   int         `json:"total_pages"`
	TotalResults int         `json:"total_results"`
}

// TMDBGenreResponse represents the response returned by the TMDB movie genre endpoint.
type TMDBGenreResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

// TMDBGenre represents a genre returned by TMDB.
type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TMDBMovie struct representing TMDb movie schema
type TMDBMovie struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	Overview         string      `json:"overview"`
	PosterPath       string      `json:"poster_path"`
	BackdropPath     string      `json:"backdrop_path"`
	ReleaseDate      string      `json:"release_date"`
	VoteAverage      float64     `json:"vote_average"`
	Runtime          int         `json:"runtime"`
	OriginalLanguage string      `json:"original_language"`
	OriginCountry    []string    `json:"origin_country"`
	Status           string      `json:"status"`
	Genres           []TMDBGenre `json:"genres"`

	Credits *Credits `json:"credits,omitempty"`
}

type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Character string `json:"character"`
}

type CrewMember struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Job  string `json:"job"`
}

// searchMovies searches TMDb by movie title.
func (s *TMDBService) SearchMovies(query string) (*TMDBSearchResponse, error) {

	encodedQuery := url.QueryEscape(query)

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/search/movie?api_key=%s&query=%s",
		s.apiKey,
		encodedQuery,
	)

	resp, err := s.get(url)
	if err != nil {
		return nil, err
	}

	// Ensure the response body is closed after we're done reading it
	defer resp.Body.Close()

	// Create a struct to hold the decoded JSON response
	var result TMDBSearchResponse

	// Decode the JSON response body into the Go struct
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// Return if the JSON is invalid or cannot be decoded
		return nil, err
	}

	// Return the populated response struct
	return &result, nil
}

// get moviedetails returns detailed  information about one movie
func (s *TMDBService) GetMovieDetails(movieID int) (*TMDBMovie, error) {

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/movie/%d?api_key=%s&append_to_response=credits",
		movieID,
		s.apiKey,
	)

	resp, err := s.get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var movie TMDBMovie

	if err := json.NewDecoder(resp.Body).Decode(&movie); err != nil {
		return nil, err
	}

	return &movie, nil
}

// PublishMovie fetches a TMDb movie and saves it to the local database if not present.
func (s *TMDBService) PublishMovie(tmdbID int) (*models.Movie, error) {

	// Step 1: Check local DB first
	existing, err := s.movieRepo.GetMovieByTMDBID(tmdbID)
	if err == nil && existing != nil {
		// If local movie has complete director and cast info, return immediately!
		if existing.Director != "" && existing.Cast != "" {
			return existing, nil
		}

		// If local movie is missing director/cast, fetch from TMDB to fill them in
		tmdbMovie, err := s.GetMovieDetails(tmdbID)
		if err == nil && tmdbMovie != nil {
			var director string
			if tmdbMovie.Credits != nil {
				for _, member := range tmdbMovie.Credits.Crew {
					if member.Job == "Director" {
						director = member.Name
						break
					}
				}
			}

			var castMembers []string
			if tmdbMovie.Credits != nil {
				for i, member := range tmdbMovie.Credits.Cast {
					if i >= 5 {
						break
					}
					castMembers = append(castMembers, member.Name)
				}
			}
			castStr := strings.Join(castMembers, ", ")

			existing.Director = director
			existing.Cast = castStr
			_ = s.movieRepo.UpdateMovie(existing)
		}

		return existing, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Step 2: If not in local DB, fetch details from TMDb
	tmdbMovie, err := s.GetMovieDetails(tmdbID)
	if err != nil {
		return nil, err
	}

	var director string
	if tmdbMovie.Credits != nil {
		for _, member := range tmdbMovie.Credits.Crew {
			if member.Job == "Director" {
				director = member.Name
				break
			}
		}
	}

	var castMembers []string
	if tmdbMovie.Credits != nil {
		for i, member := range tmdbMovie.Credits.Cast {
			if i >= 5 {
				break
			}
			castMembers = append(castMembers, member.Name)
		}
	}
	castStr := strings.Join(castMembers, ", ")

	var genres []models.Genre

	for _, genre := range tmdbMovie.Genres {
		genres = append(genres, models.Genre{
			TMDBID: genre.ID,
			Name:   genre.Name,
		})
	}

	movie := &models.Movie{
		TMDBID:       tmdbMovie.ID,
		Title:        tmdbMovie.Title,
		Description:  tmdbMovie.Overview,
		Duration:     tmdbMovie.Runtime,
		Language:     tmdbMovie.OriginalLanguage,
		ReleaseDate:  tmdbMovie.ReleaseDate,
		PosterPath:   tmdbMovie.PosterPath,
		BackdropPath: tmdbMovie.BackdropPath,
		Director:     director,
		Cast:         castStr,
		Genres:       genres,
	}

	// save movie to db
	err = s.movieRepo.CreateMovie(movie)
	if err != nil {
		return nil, err
	}

	// return the saved movie
	return movie, nil

}

// fetches official tmdb movie genres and stores any genres that are not already present in database
func (s *TMDBService) SyncGenres() error {

	url := fmt.Sprintf(
		"https://api.themoviedb.org/3/genre/movie/list?api_key=%s&language=en",
		s.apiKey,
	)

	resp, err := s.get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var result TMDBGenreResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	for _, genre := range result.Genres {

		// check if it exists locally
		_, err := s.movieRepo.GetGenreByTMDBID(genre.ID)

		if err == nil {
			continue
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// store it locally
		newGenre := &models.Genre{
			TMDBID: genre.ID,
			Name:   genre.Name,
		}

		if err := s.movieRepo.CreateGenre(newGenre); err != nil {
			return err
		}
	}

	return nil
}

// get performs a GET request with retries for transient network/API failures.
func (s *TMDBService) get(url string) (*http.Response, error) {

	const maxRetries = 3
	const baseDelay = 300 * time.Millisecond

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {

		if attempt > 0 {
			delay := baseDelay * time.Duration(1<<(attempt-1))

			// Add jitter so repeated retries are not synchronized.
			jitter := time.Duration(rand.Int63n(int64(delay)))

			time.Sleep(delay + jitter)
		}

		resp, err := s.client.Get(url)

		if err != nil {

			// Retry transient network failures such as connection resets.
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				lastErr = err
				s.client.Transport.(*http.Transport).CloseIdleConnections()
				continue
			}

			if errors.Is(err, syscall.ECONNRESET) ||
				errors.Is(err, syscall.EPIPE) ||
				errors.Is(err, io.EOF) ||
				errors.Is(err, io.ErrUnexpectedEOF) {

				lastErr = err
				s.client.Transport.(*http.Transport).CloseIdleConnections()
				continue
			}

			return nil, err
		}

		// Retry temporary server/rate-limit responses.
		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode >= 500 {

			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			lastErr = fmt.Errorf("TMDB returned status %d", resp.StatusCode)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			defer resp.Body.Close()
			return nil, fmt.Errorf("TMDB returned status %d", resp.StatusCode)
		}

		return resp, nil

	}

	return nil, fmt.Errorf(
		"TMDB request failed after %d retries: %w",
		maxRetries,
		lastErr,
	)
}
