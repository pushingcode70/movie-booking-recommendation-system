package services

import (
	"errors"
	"movie-booking/dto"
	"movie-booking/models"
	"movie-booking/repositories"
	"strings"
)

type WishlistService struct {
	repo        *repositories.WishlistRepository
	watchedRepo *repositories.WatchedMovieRepository
	tmdbService *TMDBService
}

func NewWishlistService(
	repo *repositories.WishlistRepository,
	watchedRepo *repositories.WatchedMovieRepository,
	tmdbService *TMDBService,
) *WishlistService {
	return &WishlistService{
		repo:        repo,
		watchedRepo: watchedRepo,
		tmdbService: tmdbService,
	}
}

func (s *WishlistService) AddMovie(userID uint, req *dto.CreateWishlistRequest) error {

	// already in wishlist?
	exists, err := s.repo.Exists(userID, req.TMDBID)
	if err != nil {
		return err
	}

	if exists {
		return errors.New("movie already exists in wishlist")
	}

	// check if already watched
	watched, err := s.watchedRepo.Exists(userID, req.TMDBID)
	if err != nil {
		return err
	}

	if watched {
		if err := s.watchedRepo.Delete(userID, req.TMDBID); err != nil {
			return err
		}
	}

	wishlist := &models.Wishlist{
		UserID: userID,
		TMDBID: req.TMDBID,
	}

	return s.repo.Create(wishlist)
}

func (s *WishlistService) GetWishlist(userID uint) ([]dto.WishlistItemResponse, error) {

	wishlist, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.WishlistItemResponse, 0, len(wishlist))

	for _, item := range wishlist {

		movie, err := s.tmdbService.GetMovieDetails(item.TMDBID)
		if err != nil {
			continue
		}

		movieModel := models.Movie{
			TMDBID:       movie.ID,
			Title:        movie.Title,
			Description:  movie.Overview,
			Duration:     movie.Runtime,
			Language:     movie.OriginalLanguage,
			ReleaseDate:  movie.ReleaseDate,
			PosterPath:   movie.PosterPath,
			BackdropPath: movie.BackdropPath,
			Director:     "",
			Cast:         "",
		}
		// set dynamic id to tmdb id
		movieModel.ID = uint(movie.ID)

		// parse director & cast if credits exist
		if movie.Credits != nil {
			for _, member := range movie.Credits.Crew {
				if member.Job == "Director" {
					movieModel.Director = member.Name
					break
				}
			}
			var castMembers []string
			for i, member := range movie.Credits.Cast {
				if i >= 5 {
					break
				}
				castMembers = append(castMembers, member.Name)
			}
			movieModel.Cast = strings.Join(castMembers, ", ")
		}

		response = append(response, dto.WishlistItemResponse{
			Movie:   movieModel,
			AddedAt: item.CreatedAt,
		})
	}

	return response, nil
}

func (s *WishlistService) RemoveMovie(userID uint, tmdbID int) error {
	return s.repo.Delete(userID, tmdbID)
}
