package fetcher

import (
	"fmt"
	"net/url"
	"strings"
)

type Metadata struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description,omitempty"`
	Duration    int    `json:"duration,omitempty"` // in minutes
	Type        string `json:"type"`               // video, article, etc.
	Source      string `json:"source"`             // youtube, vimeo, web, etc.
}

type Fetcher interface {
	Fetch(url string) (*Metadata, error)
	Supports(url string) bool
}

var fetchers = []Fetcher{
	&YouTubeFetcher{},
	// Add more fetchers here as needed
}

func Fetch(rawURL string) (*Metadata, error) {
	for _, f := range fetchers {
		if f.Supports(rawURL) {
			return f.Fetch(rawURL)
		}
	}
	return nil, fmt.Errorf("no fetcher available for URL: %s", rawURL)
}

func DetectSource(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "unknown"
	}

	host := strings.ToLower(parsed.Host)

	switch {
	case strings.Contains(host, "youtube.com") || strings.Contains(host, "youtu.be"):
		return "youtube"
	case strings.Contains(host, "vimeo.com"):
		return "vimeo"
	case strings.Contains(host, "spotify.com"):
		return "spotify"
	case strings.Contains(host, "goodreads.com"):
		return "goodreads"
	case strings.Contains(host, "amazon.com") || strings.Contains(host, "amazon.com.br"):
		return "amazon"
	default:
		return "web"
	}
}
