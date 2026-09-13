package services

import (
	"errors"
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
	"strings"
	"time"
)

type WatchedMovieService struct {
	repo         *repositories.WatchedMovieRepository
	wishlistRepo *repositories.WishlistRepository
	tmdbService  *TMDBService
}

func NewWatchedMovieService(
	repo *repositories.WatchedMovieRepository,
	wishlistRepo *repositories.WishlistRepository,
	tmdbService *TMDBService,
) *WatchedMovieService {
	return &WatchedMovieService{
		repo:         repo,
		wishlistRepo: wishlistRepo,
		tmdbService:  tmdbService,
	}
}

func (s *WatchedMovieService) AddWatchedMovie(userID uint, req *dto.CreateWatchedMovieRequest) error {
	exists, err := s.repo.Exists(userID, req.TMDBID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("movie already marked as watched")
	}

	watched := &models.WatchedMovie{
		UserID:    userID,
		TMDBID:    req.TMDBID,
		Rating:    req.Rating,
		Review:    req.Review,
		WatchedAt: time.Now(),
	}

	// save to watched list
	if err := s.repo.Create(watched); err != nil {
		return err
	}

	// auto-remove from wishlist
	_ = s.wishlistRepo.Delete(userID, req.TMDBID)

	return nil
}

func (s *WatchedMovieService) UpdateWatchedMovie(userID uint, tmdbID int, req *dto.UpdateWatchedMovieRequest) error {
	watched, err := s.repo.GetByUserAndTMDBID(userID, tmdbID)
	if err != nil {
		return err
	}

	watched.Rating = req.Rating
	watched.Review = req.Review

	return s.repo.Update(watched)
}

func (s *WatchedMovieService) GetWatchedMovies(userID uint) ([]dto.WatchedMovieItemResponse, error) {
	items, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	var response []dto.WatchedMovieItemResponse
	for _, item := range items {
		tmdbMovie, err := s.tmdbService.GetMovieDetails(item.TMDBID)
		if err != nil {
			continue
		}

		movieModel := models.Movie{
			TMDBID:       tmdbMovie.ID,
			Title:        tmdbMovie.Title,
			Description:  tmdbMovie.Overview,
			Duration:     tmdbMovie.Runtime,
			Language:     tmdbMovie.OriginalLanguage,
			ReleaseDate:  tmdbMovie.ReleaseDate,
			PosterPath:   tmdbMovie.PosterPath,
			BackdropPath: tmdbMovie.BackdropPath,
			Director:     "",
			Cast:         "",
		}
		movieModel.ID = uint(tmdbMovie.ID)

		if tmdbMovie.Credits != nil {
			for _, member := range tmdbMovie.Credits.Crew {
				if member.Job == "Director" {
					movieModel.Director = member.Name
					break
				}
			}
			var castMembers []string
			for i, member := range tmdbMovie.Credits.Cast {
				if i >= 5 {
					break
				}
				castMembers = append(castMembers, member.Name)
			}
			movieModel.Cast = strings.Join(castMembers, ", ")
		}

		response = append(response, dto.WatchedMovieItemResponse{
			Movie:     movieModel,
			Rating:    item.Rating,
			Review:    item.Review,
			WatchedAt: item.WatchedAt,
		})
	}

	return response, nil
}

func (s *WatchedMovieService) RemoveWatchedMovie(userID uint, tmdbID int) error {
	return s.repo.Delete(userID, tmdbID)
}
