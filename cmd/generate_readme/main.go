package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

type FormattedItem struct {
	Title  string
	Author *string
	Tags   string
	URL    string
}

var translations = map[string]string{
	// Type
	"article": "Artigo",
	"book":    "Livro",
	"course":  "Curso",
	"feed":    "Feed",
	"podcast": "Podcast",
	"roadmap": "Roadmap",
	"video":   "Vídeo",

	// Level
	"beginner":     "Iniciante",
	"intermediate": "Intermediário",
	"advanced":     "Avançado",
}

func main() {
	data, err := os.ReadFile("data.yml")
	if err != nil {
		fmt.Println("error reading data.yml:", err)
		os.Exit(1)
	}

	items, err := catalog.ParseCatalog(data)
	if err != nil {
		fmt.Println("error parsing catalog:", err)
		os.Exit(1)
	}

	readme, err := GenerateREADME(items)
	if err != nil {
		fmt.Println("error generating readme:", err)
		os.Exit(1)
	}

	fmt.Println("readme:", readme)

	fmt.Print("write readme.md file")

	err = os.WriteFile("README.md", []byte(readme), 0644)
	if err != nil {
		fmt.Println("error writing README.md", err)
		os.Exit(1)
	}
}

func translate(key string) string {
	if val, ok := translations[key]; ok {
		return val
	}

	return key
}

type FilterItemFn func(catalog.CatalogItem) bool

func filterItem(items []catalog.CatalogItem, predicate FilterItemFn) []catalog.CatalogItem {
	filtered := make([]catalog.CatalogItem, 0)

	for _, item := range items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

func removeTag(item catalog.CatalogItem, tag string) catalog.CatalogItem {
	tagInLowerCase := strings.ToLower(tag)
	newTags := make([]string, 0, len(item.Tags))

	for _, t := range item.Tags {
		if strings.ToLower(t) != tagInLowerCase {
			newTags = append(newTags, t)
		}
	}

	return catalog.CatalogItem{
		Title:  item.Title,
		URL:    item.URL,
		Type:   item.Type,
		Tags:   newTags,
		Level:  item.Level,
		IsPaid: item.IsPaid,
		Author: item.Author,
	}
}

func hasTag(item catalog.CatalogItem, tag string) bool {
	tagInLowerCase := strings.ToLower(tag)

	for _, v := range item.Tags {
		if strings.ToLower(v) == tagInLowerCase {
			return true
		}
	}

	return false
}

func formatCatalogItems(items []catalog.CatalogItem) []FormattedItem {
	var formattedItems []FormattedItem

	for _, item := range items {
		item = removeTag(item, "Excelência Técnica")
		item = removeTag(item, "Liderança e Inspiração")
		item = removeTag(item, "Entrega de Valor")

		item.Tags = append(item.Tags, translate(item.Level))
		item.Tags = append(item.Tags, translate(item.Type))

		if item.IsPaid {
			item.Tags = append(item.Tags, "Pago")
		} else {
			item.Tags = append(item.Tags, "Grátis")
		}

		formattedItems = append(formattedItems, FormattedItem{
			Title:  getTitle(item),
			Author: item.Author,
			Tags:   formatTags(item.Tags),
			URL:    item.URL,
		})
	}

	return formattedItems
}

func getTitle(item catalog.CatalogItem) string {
	title := strings.ReplaceAll(item.Title, "|", "-")
	title = strings.ReplaceAll(title, "\n", " ")

	return strings.TrimSpace(title)
}

func formatTags(tags []string) string {
	newTags := make([]string, 0, len(tags))

	for _, t := range tags {
		newTags = append(newTags, fmt.Sprintf("`%s`", t))
	}

	return safeJoin(newTags, " ")
}

func safeJoin(slice []string, sep string) string {
	if slice == nil {
		return ""
	}

	return strings.Join(slice, sep)
}

func GenerateREADME(items []catalog.CatalogItem) (string, error) {
	sort.Slice(items, func(i, j int) bool {
		return strings.ToLower(items[i].Title) < strings.ToLower(items[j].Title)
	})

	technicalExcellence := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return hasTag(item, "Excelência Técnica")
	}))

	leadershipAndInspiration := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return hasTag(item, "Liderança e Inspiração")
	}))

	deliveringValue := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return hasTag(item, "Entrega de Valor")
	}))

	const readmeTemplate = `
# Awesome Tech Lead [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)

Uma lista de conteúdo sobre lideraça técnica curada pelos membros da comunidade [TechLeads.club 💎](https://comece.techleads.club?utm_source=awesome-tech-lead&utm_medium=readme).

{{if .TechnicalExcellence}}
## 🏆 Excelência Técnica

| Título                                                          | Tags        | 
|-----------------------------------------------------------------|-------------|
{{- range .TechnicalExcellence }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}

{{if .DeliveringValue}}
## 📦 Entrega de Valor 

| Título                                                          | Tags        |
|-----------------------------------------------------------------|-------------|
{{- range .DeliveringValue }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}

{{if .LeadershipAndInspiration}}
## 🤝 Liderança e Inspiração 

| Título                                                          | Tags        |
|-----------------------------------------------------------------|-------------|
{{- range .LeadershipAndInspiration }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	templateData := map[string][]FormattedItem{
		"TechnicalExcellence":      technicalExcellence,
		"DeliveringValue":          deliveringValue,
		"LeadershipAndInspiration": leadershipAndInspiration,
	}

	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}
