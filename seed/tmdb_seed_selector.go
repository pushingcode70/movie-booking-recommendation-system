package seed

import (
	"sort"
)

func SelectTopMoviesByPopularity(movies []TMDBExportMovie, limit int) []TMDBExportMovie {

	if limit <= 0 {
		return []TMDBExportMovie{}
	}

	if limit > len(movies) {
		limit = len(movies)
	}

	sort.Slice(movies, func(i, j int) bool {
		return movies[i].Popularity > movies[j].Popularity
	})

	return movies[:limit]
}
