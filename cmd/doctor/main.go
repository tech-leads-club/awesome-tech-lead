package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

const (
	DefaultTimeout     = 15 * time.Second
	DefaultConcurrency = 50
	UserAgent          = "AwesomeTechLead-Doctor/1.0 (+https://github.com/tech-leads-club/awesome-tech-lead)"
)

type CheckStatus int

const (
	StatusOK CheckStatus = iota
	StatusBroken
	StatusTimeout
	StatusDNSError
	StatusConnectionError
	StatusTLSError
)

func (s CheckStatus) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusBroken:
		return "BROKEN"
	case StatusTimeout:
		return "TIMEOUT"
	case StatusDNSError:
		return "DNS_ERROR"
	case StatusConnectionError:
		return "CONNECTION_ERROR"
	case StatusTLSError:
		return "TLS_ERROR"
	default:
		return "UNKNOWN"
	}
}

type URLCheckResult struct {
	URL        string
	Title      string
	StatusCode int
	Status     CheckStatus
	Error      error
	Duration   time.Duration
}

type URLChecker struct {
	client      *http.Client
	concurrency int
}

func NewURLChecker(timeout time.Duration, concurrency int) *URLChecker {
	return &URLChecker{
		client: &http.Client{
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 10 {
					return fmt.Errorf("too many redirects")
				}
				return nil
			},
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		concurrency: concurrency,
	}
}

func (c *URLChecker) CheckURL(ctx context.Context, url, title string) URLCheckResult {
	start := time.Now()
	result := URLCheckResult{
		URL:   url,
		Title: title,
	}

	// Try HEAD first, fallback to GET if server rejects HEAD
	result = c.doRequest(ctx, http.MethodHead, url, title)
	result.Duration = time.Since(start)

	// Fallback to GET for sites that reject HEAD requests or timeout
	if result.StatusCode == 403 || result.StatusCode == 405 || result.StatusCode == 406 || result.Status == StatusTimeout {
		start = time.Now()
		result = c.doRequest(ctx, http.MethodGet, url, title)
		result.Duration = time.Since(start)
	}

	return result
}

func (c *URLChecker) doRequest(ctx context.Context, method, url, title string) URLCheckResult {
	result := URLCheckResult{
		URL:   url,
		Title: title,
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		result.Status = StatusConnectionError
		result.Error = err
		return result
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := c.client.Do(req)
	if err != nil {
		result.Error = err
		result.Status = classifyError(err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 400 {
		result.Status = StatusBroken
	} else {
		result.Status = StatusOK
	}

	return result
}

func classifyError(err error) CheckStatus {
	errStr := strings.ToLower(err.Error())
	switch {
	case strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded"):
		return StatusTimeout
	case strings.Contains(errStr, "no such host") || strings.Contains(errStr, "dns"):
		return StatusDNSError
	case strings.Contains(errStr, "certificate") || strings.Contains(errStr, "tls") || strings.Contains(errStr, "x509"):
		return StatusTLSError
	default:
		return StatusConnectionError
	}
}

type ProgressFunc func(checked, total int, result URLCheckResult)

func (c *URLChecker) CheckAll(ctx context.Context, items []catalog.CatalogItem, onProgress ProgressFunc) []URLCheckResult {
	results := make([]URLCheckResult, len(items))
	total := len(items)

	sem := make(chan struct{}, c.concurrency)
	var wg sync.WaitGroup
	var checked int
	var mu sync.Mutex

	for i, item := range items {
		wg.Add(1)
		go func(idx int, item catalog.CatalogItem) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			result := c.CheckURL(ctx, item.URL, item.Title)
			results[idx] = result

			if onProgress != nil {
				mu.Lock()
				checked++
				onProgress(checked, total, result)
				mu.Unlock()
			}
		}(i, item)
	}

	wg.Wait()
	return results
}

type Report struct {
	TotalChecked int
	BrokenCount  int
	OKCount      int
	Results      []URLCheckResult
}

func generateReport(results []URLCheckResult) Report {
	report := Report{
		TotalChecked: len(results),
		Results:      results,
	}

	for _, r := range results {
		if r.Status == StatusOK {
			report.OKCount++
		} else {
			report.BrokenCount++
		}
	}

	return report
}

func printReport(report Report) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("CATALOG HEALTH CHECK REPORT")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println()

	fmt.Printf("Total URLs checked: %d\n", report.TotalChecked)
	fmt.Printf("Working:            %d\n", report.OKCount)
	fmt.Printf("Broken:             %d\n", report.BrokenCount)
	fmt.Println()

	if report.BrokenCount > 0 {
		fmt.Println(strings.Repeat("-", 80))
		fmt.Println("BROKEN LINKS:")
		fmt.Println(strings.Repeat("-", 80))

		for _, r := range report.Results {
			if r.Status != StatusOK {
				fmt.Printf("\n[%s] %s\n", r.Status, r.Title)
				fmt.Printf("  URL:    %s\n", r.URL)
				if r.StatusCode > 0 {
					fmt.Printf("  Status: %d\n", r.StatusCode)
				}
				if r.Error != nil {
					fmt.Printf("  Error:  %s\n", r.Error)
				}
			}
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("=", 80))
}

func main() {
	var concurrency int
	var showProgress bool

	flag.IntVar(&concurrency, "concurrency", DefaultConcurrency, "number of concurrent requests")
	flag.IntVar(&concurrency, "c", DefaultConcurrency, "number of concurrent requests (shorthand)")
	flag.BoolVar(&showProgress, "progress", false, "show progress as URLs are checked")
	flag.Parse()

	data, err := os.ReadFile("catalog.yml")
	if err != nil {
		fmt.Println("error reading catalog.yml:", err)
		os.Exit(1)
	}

	items, err := catalog.ParseCatalog(data)
	if err != nil {
		fmt.Println("error parsing catalog:", err)
		os.Exit(1)
	}

	fmt.Printf("Checking %d URLs in catalog (concurrency: %d)...\n", len(items), concurrency)

	checker := NewURLChecker(DefaultTimeout, concurrency)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var progressFn ProgressFunc
	if showProgress {
		progressFn = func(checked, total int, result URLCheckResult) {
			status := "OK"
			if result.Status != StatusOK {
				status = result.Status.String()
			}
			fmt.Printf("\r[%d/%d] %s - %s", checked, total, status, truncate(result.Title, 50))
			// Clear rest of line
			fmt.Print("\033[K")
		}
	}

	results := checker.CheckAll(ctx, items, progressFn)

	if showProgress {
		fmt.Println() // New line after progress
	}

	report := generateReport(results)
	printReport(report)

	if report.BrokenCount > 0 {
		os.Exit(1)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
