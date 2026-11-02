package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type flags struct {
	from    int
	to      int
	workers int
	timeout time.Duration
}

func main() {

	fromPtr := flag.Int("from", -1, "mandatory index flag of the first movie")
	toPtr := flag.Int("to", -1, "mandatory index flag of the last movie")
	workersPtr := flag.Int("workers", 10, "number of workers")
	timeoutPtr := flag.Duration("timeout", 5*time.Second, "timeout of a single HTTP request")

	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		seen[(*f).Name] = true
	})
	if !seen["from"] {
		fmt.Fprintln(os.Stderr, "the required value of the \"from\" flag is not specified")
		os.Exit(2)
	}
	if !seen["to"] {
		fmt.Fprintln(os.Stderr, "the required value of the \"to\" flag is not specified")
		os.Exit(2)
	}

	if (*fromPtr) < 1 ||
		(*toPtr) < 1 ||
		(*workersPtr) < 1 {
		fmt.Fprintln(os.Stderr, "the value of the \"from\", \"to\", \"workers\" flags must be at least 1")
		os.Exit(2)
	}

	if *timeoutPtr <= 0 {
		fmt.Fprintf(os.Stderr, "timeout must be > 0, got %s\n", *timeoutPtr)
		os.Exit(2)
	}

	if (*fromPtr) > (*toPtr) {
		fmt.Fprintln(os.Stderr, "the value of the \"from\" parameter is always less than or equal to the value of the \"to\" parameter")
		os.Exit(2)
	}

	cfg := flags{
		from:    *fromPtr,
		to:      *toPtr,
		workers: *workersPtr,
		timeout: *timeoutPtr,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	code := run(ctx, cfg)

	stop()
	os.Exit(code)
}
