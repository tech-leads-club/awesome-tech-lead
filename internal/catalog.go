package catalog

import (
	"fmt"

	"gopkg.in/yaml.v3"
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

	// Parse YAML data
	if err := yaml.Unmarshal(data, &catalog); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	// Validate each item
	for _, item := range catalog.Catalog {
		if err := validateCatalogItem(item); err != nil {
			return nil, fmt.Errorf("validation error for item %q: %w", item.Title, err)
		}
	}

	return catalog.Catalog, nil
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

	return nil
}
