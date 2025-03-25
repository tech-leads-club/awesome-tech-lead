package catalog

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/gosimple/slug"
)

const (
	TechnicalExcellenceTag      = "Excelência Técnica"
	LeadershipAndInspirationTag = "Liderança e Inspiração"
	DeliveringValueTag          = "Entrega de Valor"
)

var PillarTags = map[string]bool{
	TechnicalExcellenceTag:      true,
	LeadershipAndInspirationTag: true,
	DeliveringValueTag:          true,
}

var ValidTypes = map[string]struct{}{
	"article": {},
	"book":    {},
	"course":  {},
	"feed":    {},
	"podcast": {},
	"roadmap": {},
	"video":   {},
}

var ValidCareerBands = map[string]struct{}{
	"junior":    {},
	"mid":       {},
	"senior":    {},
	"tl":        {},
	"staff":     {},
	"principal": {},
}

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

		if seenURLs[item.URL] {
			return nil, fmt.Errorf("duplicate URL found: %s", item.URL)
		}
		seenURLs[item.URL] = true

		slugTitle := slug.Make(item.Title)
		if seenTitles[slugTitle] {
			return nil, fmt.Errorf("duplicate title found: %s", item.Title)
		}
		seenTitles[slugTitle] = true
	}

	return catalog.Catalog, nil
}

func validateCatalogItem(item CatalogItem) error {
	var errs []string

	if _, ok := ValidTypes[item.Type]; !ok {
		errs = append(errs, fmt.Sprintf("invalid type %q, valid options are: %s", item.Type, joinKeys(ValidTypes)))
	}

	if len(item.CareerBands) == 0 {
		errs = append(errs, "career bands cannot be empty")
	} else {
		for _, band := range item.CareerBands {
			if _, ok := ValidCareerBands[band]; !ok {
				errs = append(errs, fmt.Sprintf("invalid career band %q, valid options are: %s", band, joinKeys(ValidCareerBands)))
			}
		}
	}

	if len(item.Tags) == 0 {
		errs = append(errs, "tags cannot be empty")
	} else if !hasPillarTag(item.Tags) {
		errs = append(errs, fmt.Sprintf("item must have at least one pillar tag: %s", joinPillarTags()))
	}

	if len(errs) > 0 {
		return fmt.Errorf("validation errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

func hasPillarTag(tags []string) bool {
	for _, tag := range tags {
		if PillarTags[tag] {
			return true
		}
	}
	return false
}

func joinKeys(m map[string]struct{}) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func joinPillarTags() string {
	return fmt.Sprintf("%q, %q, %q", TechnicalExcellenceTag, LeadershipAndInspirationTag, DeliveringValueTag)
}
