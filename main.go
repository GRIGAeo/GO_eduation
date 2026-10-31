package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Movie struct {
	id       int
	title    string
	year     int
	director string
}

func main() {

}

func fetchMovie(ctx context.Context, id int32, client http.Client) (Movie, error) {
	url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Movie{}, fmt.Errorf("movie: %d, error: %w", id, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return Movie{}, fmt.Errorf("movie: %d, error: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Movie{}, fmt.Errorf("movie: %d, HTTP STATUS: %d", id, resp.StatusCode)
	}

	var m = Movie{}
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil{
		return Movie{}, fmt.Errorf("movie: %d, invalid JSON: %d", err)
	}

	return m, nil
}
