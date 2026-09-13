package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"movie-booking/config"
	"movie-booking/database"
	"movie-booking/models"
	"movie-booking/repositories"
	"movie-booking/seed"
	"movie-booking/services"

	"github.com/joho/godotenv"
	"github.com/pgvector/pgvector-go"
)

const (
	workerCount = 4
	movieLimit  = 20000 // number of movies sorted top in file (20k)
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("usage: go run ./cmd/seed <path-to-tmdb-export.json.gz>")
	}

	filePath := os.Args[1]

	// read tmdb daily export
	movies, err := seed.ReadMovieExport(filePath)
	if err != nil {
		log.Fatal("failed to read TMDB export: ", err)
	}

	fmt.Println("Total records:", len(movies))

	config.LoadConfig()

	// select most popular movies for recommendation catalogue
	seedMovies := seed.SelectTopMoviesByPopularity(movies, movieLimit)

	fmt.Println("Selected seed movies:", len(seedMovies))
	fmt.Println("Workers:", workerCount)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()

	movieRepo := repositories.NewMovieRepository(database.DB)
	tmdbService := services.NewTMDBService(movieRepo)

	embeddingService := services.NewMovieEmbeddingService(tmdbService)
	embeddingRepo := repositories.NewMovieEmbeddingRepository(database.DB)

	jobs := make(chan int)

	var wg sync.WaitGroup
	var mu sync.Mutex

	processed := 0
	skipped := 0
	failed := 0

	for workerID := 1; workerID <= workerCount; workerID++ {

		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for tmdbID := range jobs {

				// skip if movie and embedding already exist
				existingMovie, movieErr := movieRepo.GetMovieByTMDBID(tmdbID)

				_, embeddingErr := embeddingRepo.GetByTMDBID(tmdbID)

				if movieErr == nil && embeddingErr == nil {
					mu.Lock()
					skipped++
					mu.Unlock()

					fmt.Printf(
						"Worker %d: SKIPPED TMDB ID=%d (%s)\n",
						id,
						tmdbID,
						existingMovie.Title,
					)

					continue
				}

				// get movie metadata, genres and credits in one tmdb request
				tmdbMovie, err := tmdbService.GetMovieDetails(tmdbID)
				if err != nil {
					mu.Lock()
					failed++
					mu.Unlock()

					fmt.Printf(
						"Worker %d: FAILED TMDB ID=%d - %v\n",
						id,
						tmdbID,
						err,
					)

					continue
				}

				// extract director
				var director string

				if tmdbMovie.Credits != nil {
					for _, member := range tmdbMovie.Credits.Crew {
						if member.Job == "Director" {
							director = member.Name
							break
						}
					}
				}

				// extract up to five cast members
				var castMembers []string

				if tmdbMovie.Credits != nil {
					for i, member := range tmdbMovie.Credits.Cast {
						if i >= 5 {
							break
						}

						castMembers = append(castMembers, member.Name)
					}
				}

				cast := strings.Join(castMembers, ", ")

				// build local movie model
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
					Cast:         cast,
				}

				for _, genre := range tmdbMovie.Genres {
					movie.Genres = append(movie.Genres, models.Genre{
						TMDBID: genre.ID,
						Name:   genre.Name,
					})
				}

				// create movie if it does not exist
				if movieErr != nil {
					if err := movieRepo.CreateMovie(movie); err != nil {
						mu.Lock()
						failed++
						mu.Unlock()

						fmt.Printf(
							"Worker %d: FAILED MOVIE TMDB ID=%d - %v\n",
							id,
							tmdbID,
							err,
						)

						continue
					}
				}

				// generate embedding from tmdb response
				text := embeddingService.BuildMovieText(tmdbMovie)

				embedding, err := embeddingService.GenerateEmbedding(text)
				if err != nil {
					mu.Lock()
					failed++
					mu.Unlock()

					fmt.Printf(
						"Worker %d: FAILED EMBEDDING TMDB ID=%d - %v\n",
						id,
						tmdbID,
						err,
					)

					continue
				}

				// create embedding if it does not exist
				if embeddingErr != nil {
					movieEmbedding := &models.MovieEmbedding{
						TMDBID:    tmdbID,
						Embedding: pgvector.NewVector(embedding),
					}

					if err := embeddingRepo.Create(movieEmbedding); err != nil {
						mu.Lock()
						failed++
						mu.Unlock()

						fmt.Printf(
							"Worker %d: FAILED EMBEDDING SAVE TMDB ID=%d - %v\n",
							id,
							tmdbID,
							err,
						)

						continue
					}
				}

				mu.Lock()
				processed++
				mu.Unlock()

				fmt.Printf(
					"Worker %d: COMPLETED TMDB ID=%d (%s)\n",
					id,
					tmdbID,
					tmdbMovie.Title,
				)
			}
		}(workerID)
	}

	for _, movie := range seedMovies {
		jobs <- movie.ID
	}

	close(jobs)

	wg.Wait()

	fmt.Println()
	fmt.Println("TMDB catalogue seeding completed.")
	fmt.Println("Processed:", processed)
	fmt.Println("Skipped:", skipped)
	fmt.Println("Failed:", failed)
}
