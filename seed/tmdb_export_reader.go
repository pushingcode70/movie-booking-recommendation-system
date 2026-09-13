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

// reads TMDB's gzipped daily movie export.
func ReadMovieExport(filePath string) ([]TMDBExportMovie, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decompress the .json.gz export.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	var movies []TMDBExportMovie

	scanner := bufio.NewScanner(gzReader)

	//increase the limit because export records can be large.
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	//each line in the export is a separate JSON object.
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
