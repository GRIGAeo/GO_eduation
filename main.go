package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Movie struct {
	id       int
	title    string
	year     int
	director string
}

func (mov *Movie) mOutput(){
	fmt.Println("%d — %s — %d — %s", mov.id, mov.title, mov.year, mov.director)
}

func main() int{

	fromPtr  := flag.Int("from", -1, "mandatory index flag of the first movie")
	toPtr := flag.Int("to", -1, "mandatory index flag of the last movie")
	workersPtr := flag.Int("workers", 10, "number of workers")
	timeout := flag.Duration("timeout", 5*time.Second, "timeout of a single HTTP request")

	flag.Parse()
	
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag){
		seen[(*f).Name] = true
	})
	if !seen["from"] && !seen["to"]{
		fmt.Println(os.Stderr, "the required values of the \"from\" and \"to\" flags are not specified")
		return 2
	}
	if !seen["from"]{
		fmt.Println(os.Stderr, "the required value of the \"from\" flag is not specified")
		return 2
	}
	if !seen["to"]{
		fmt.Println(os.Stderr, "the required value of the \"to\" flag is not specified")
		return 2
	}

	if (*fromPtr) < 1|| 
	(*toPtr) < 1 ||
	(*workersPtr) < 1{
		fmt.Println(os.Stderr, "the value of the \"from\", \"to\", \"workers\" flags must be at least 1")
		return 2
	}

	if(*fromPtr) > (*toPtr){
		fmt.Println(os.Stderr, "the value of the \"from\" parameter is always less than or equal to the value of the \"to\" parameter")
		return 2
	}


}

func fetchMovie(ctx context.Context, id int, client http.Client) (Movie, error) {
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
