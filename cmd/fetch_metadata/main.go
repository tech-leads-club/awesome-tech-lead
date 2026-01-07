package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/tech-leads-club/awesome-tech-lead/internal/fetcher"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: fetch_metadata <url> [url...]")
		fmt.Fprintln(os.Stderr, "\nFetches metadata from URLs and outputs JSON.")
		fmt.Fprintln(os.Stderr, "\nSupported sources:")
		fmt.Fprintln(os.Stderr, "  - YouTube (youtube.com, youtu.be)")
		fmt.Fprintln(os.Stderr, "\nExamples:")
		fmt.Fprintln(os.Stderr, "  fetch_metadata https://www.youtube.com/watch?v=xyz")
		fmt.Fprintln(os.Stderr, "  fetch_metadata https://youtu.be/xyz https://youtu.be/abc")
		os.Exit(1)
	}

	urls := os.Args[1:]

	if len(urls) == 1 {
		// Single URL: output single object
		metadata, err := fetcher.Fetch(urls[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		outputJSON(metadata)
	} else {
		// Multiple URLs: output array
		var results []*fetcher.Metadata
		for _, url := range urls {
			metadata, err := fetcher.Fetch(url)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching %s: %v\n", url, err)
				continue
			}
			results = append(results, metadata)
		}
		outputJSON(results)
	}
}

func outputJSON(v interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(v); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}
