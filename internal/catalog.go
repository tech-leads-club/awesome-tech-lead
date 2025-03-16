package catalog

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"net/url"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"gopkg.in/yaml.v3"
)

const (
        TechnicalExcellenceTag      = "Excelência Técnica"
        LeadershipAndInspirationTag = "Liderança e Inspiração"
        DeliveringValueTag          = "Entrega de Valor"
)

type CatalogItem struct {
	URL         string   `yaml:"url"`
	Title       string   `yaml:"title"`
	Author      *string  `yaml:"author,omitempty"`
	Type        string   `yaml:"type"`
	Tags        []string `yaml:"tags"`
	IsPaid      bool     `yaml:"is_paid"`
	Level       string   `yaml:"level"`
	CareerBands []string `yaml:"career_bands"`
	Language    string   `yaml:"language"`
	Duration    *string  `yaml:"duration,omitempty"`
}

func ParseCatalog(data []byte) ([]CatalogItem, error) {
	var catalog struct {
		Catalog []CatalogItem `yaml:"catalog"`
	}

	seenURLs := make(map[string]bool)
	seenTitles := make(map[string]bool)

	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	for _, item := range catalog.Catalog {
		if err := validateCatalogItem(item); err != nil {
			return nil, fmt.Errorf("validation error for item %q: %w", item.Title, err)
		}

		cleanedURL, err := cleanURL(item.URL)
		if err != nil {
			fmt.Println("validation error URL:", err)
			continue
		}

		if seenURLs[cleanedURL] {
			return nil, fmt.Errorf("duplicate URL found: %s", item.URL)
		}
		seenURLs[cleanedURL] = true

		slugTitle := generateSlug(item.Title)
		if seenTitles[slugTitle] {
			return nil, fmt.Errorf("duplicate title found: %s", item.Title)
		}
		seenTitles[slugTitle] = true
	}

	return catalog.Catalog, nil
}

func cleanURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""

	return strings.TrimSuffix(parsedURL.String(), "/"), nil
}

func generateSlug(title string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, title)

	result = strings.ToLower(result)

	reg, _ := regexp.Compile("[^a-z0-9-]+")
	result = reg.ReplaceAllString(result, "-")

	result = strings.Trim(regexp.MustCompile(`-+`).ReplaceAllString(result, "-"), "-")

	return result
}

func validateCatalogItem(item CatalogItem) error {
	allowedTypes := map[string]bool{
		"article": true,
		"book":    true,
		"course":  true,
		"feed":    true,
		"podcast": true,
		"roadmap": true,
		"video":   true,
	}

	if !allowedTypes[item.Type] {
		return fmt.Errorf("invalid type: %s", item.Type)
	}

	if len(item.Tags) == 0 {
		return fmt.Errorf("tags cannot be empty")
	}

	pillarTags := map[string]bool{
		TechnicalExcellenceTag:      true,
		LeadershipAndInspirationTag: true,
		DeliveringValueTag:          true,
	}

	for _, tag := range item.Tags {
		if pillarTags[tag] {
			return nil
		}
	}

	return fmt.Errorf("item must have at least one pillar tag: %s, %s, or %s", 
		TechnicalExcellenceTag, LeadershipAndInspirationTag, DeliveringValueTag)
}
