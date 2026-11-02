package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type result struct {
	id  int
	mv  movie
	err error
}

func run(ctx context.Context, flgs flags) int {
	client := &http.Client{Timeout: flgs.timeout}

	jobs := make(chan int)
	results := make(chan result)

	var wg sync.WaitGroup

	for i := 0; i < flgs.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				m, err := fetchMovie(ctx, id, client)
				select {
				case results <- result{id: id, mv: m, err: err}:
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for id := flgs.from; id <= flgs.to; id++ {
			select {
			case jobs <- id:
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if r.err != nil {
			if errors.Is(r.err, context.Canceled) {
				continue
			}
			fmt.Fprintf(os.Stderr, "%v\n", r.err)
			continue
		}
		r.mv.mOutput()
	}

	if ctx.Err() != nil {
		fmt.Fprintln(os.Stderr, "interrupted")
		return 130
	}

	return 0
}
