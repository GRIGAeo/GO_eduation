package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type movie struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func (mov *movie) mOutput() {
	fmt.Printf("%d — %s — %d — %s\n", mov.ID, mov.Title, mov.Year, mov.Director)
}

func fetchMovie(ctx context.Context, id int, client *http.Client) (movie, error) {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return movie{}, fmt.Errorf("movie: %d, error: %w", id, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return movie{}, fmt.Errorf("movie: %d, error: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return movie{}, fmt.Errorf("movie: %d, HTTP STATUS: %d", id, resp.StatusCode)
	}

	var m = movie{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return movie{}, fmt.Errorf("movie: %d, invalid JSON: %w", id, err)
	}

	return m, nil
}
