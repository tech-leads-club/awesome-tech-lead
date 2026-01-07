package fetcher

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"strings"
)

type YouTubeFetcher struct{}

type ytdlpOutput struct {
	Title       string `json:"title"`
	Channel     string `json:"channel"`
	Uploader    string `json:"uploader"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	WebpageURL  string `json:"webpage_url"`
}

func (f *YouTubeFetcher) Supports(rawURL string) bool {
	return DetectSource(rawURL) == "youtube"
}

func (f *YouTubeFetcher) Fetch(rawURL string) (*Metadata, error) {
	cmd := exec.Command("yt-dlp",
		"-j",
		"--no-download",
		"--no-playlist",
		rawURL,
	)

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("yt-dlp failed: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("yt-dlp execution error: %w", err)
	}

	var yt ytdlpOutput
	if err := json.Unmarshal(output, &yt); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}

	title, author := extractTitleAndAuthor(yt.Title, yt.Channel, yt.Uploader)

	durationMinutes := int(math.Ceil(float64(yt.Duration) / 60))

	return &Metadata{
		URL:         yt.WebpageURL,
		Title:       title,
		Author:      author,
		Description: yt.Description,
		Duration:    durationMinutes,
		Type:        "video",
		Source:      "youtube",
	}, nil
}

func extractTitleAndAuthor(title, channel, uploader string) (string, string) {
	// Try to extract speaker name from title patterns like:
	// "Title - Speaker Name" or "Title | Speaker Name"
	for _, sep := range []string{" - ", " | ", " – "} {
		if idx := strings.LastIndex(title, sep); idx != -1 {
			candidate := strings.TrimSpace(title[idx+len(sep):])
			// Check if it looks like a name (not too long, no special chars)
			if len(candidate) < 50 && !strings.ContainsAny(candidate, "()[]{}") {
				cleanTitle := strings.TrimSpace(title[:idx])
				return cleanTitle, candidate
			}
		}
	}

	// Fallback: use full title and channel/uploader as author
	author := channel
	if author == "" {
		author = uploader
	}

	return title, author
}
