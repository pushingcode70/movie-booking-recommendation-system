package seed

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
)

type TMDBExportMovie struct {
	ID         int     `json:"id"`
	Popularity float64 `json:"popularity"`
}

// reads tmdb's gzipped daily movie export
func ReadMovieExport(filePath string) ([]TMDBExportMovie, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decompress the .json.gz export
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	var movies []TMDBExportMovie

	scanner := bufio.NewScanner(gzReader)

	// increase buffer limit for large export records
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	// each line in export is a separate json object
	for scanner.Scan() {

		var movie TMDBExportMovie

		if err := json.Unmarshal(scanner.Bytes(), &movie); err != nil {
			return nil, fmt.Errorf("failed to parse export record: %w", err)
		}

		movies = append(movies, movie)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
