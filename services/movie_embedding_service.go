package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"movie-booking/config"
)

type MovieEmbeddingService struct {
	client      *http.Client
	embedURL    string
	tmdbService *TMDBService
}

func NewMovieEmbeddingService(tmdbService *TMDBService) *MovieEmbeddingService {
	return &MovieEmbeddingService{
		client:      &http.Client{},
		embedURL:    config.AppConfig.EmbeddingServiceURL,
		tmdbService: tmdbService,
	}
}

func (s *MovieEmbeddingService) BuildMovieText(movie *TMDBMovie) string {

	var genreNames []string

	for _, genre := range movie.Genres {
		genreNames = append(genreNames, genre.Name)
	}

	return fmt.Sprintf(

		"Overview: %s\n"+
			"Genres: %s",

		movie.Overview,
		strings.Join(genreNames, ", "),
	)
}

func (s *MovieEmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	//go struct -> json.Marshal() -> json -> python ->python object ->model.encode() ->embedding vector -> json response -> json.NewDecoder().Decode() -> go struct -> []float32
	//build josn request body expected by embedding service
	requestBody := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}

	//converts go struct into json bytes
	data, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Post(
		s.embedURL,            //where to send ..this is url of python endpoint
		"application/json",    //format.. header and type
		bytes.NewBuffer(data), //wraps json bytes	`data` in buffer..http client needs the requst body as a readable stream,
	// so the buffer allows those bytes  to be used as http body
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // if we dont use defer we have close after each return type defer says run this when function is about to end

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embedding service returned status %d", resp.StatusCode)
	}

	//decode the json response containing the embedding vector
	var result struct {
		Embedding []float32 `json:"embedding"`
	}

	//resp.body contains the body of  of http response 	recieved from python
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil { //read the json body and converts it into go struct
		return nil, err
	}

	return result.Embedding, nil //extracts the embedding field from go struct  which is now represented as go slice containing floating point numbers

}

func (s *MovieEmbeddingService) CreateMovieEmbedding(movie *TMDBMovie) ([]float32, error) {

	text := s.BuildMovieText(movie)
	return s.GenerateEmbedding(text)
}

func (s *MovieEmbeddingService) CreateMovieEmbeddingByTMDBID(tmdbID int) ([]float32, error) {

	movie, err := s.tmdbService.GetMovieDetails(tmdbID)
	if err != nil {
		return nil, err
	}

	return s.CreateMovieEmbedding(movie)
}
