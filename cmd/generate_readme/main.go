package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

type FormattedItem struct {
	Title  string
	Author *string
	Type   string
	Tags   string
	Level  string
	IsPaid string
	URL    string
}

var translations = map[string]string{
	// Type
	"article": "Artigo",
	"book":    "Livro",
	"course":  "Curso",
	"feed":    "Feed",
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

func GenerateREADME(items []catalog.CatalogItem) (string, error) {
	books := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "book"
	}))

	articles := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "article"
	}))

	courses := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "course"
	}))

	videos := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "video"
	}))

	podcasts := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "podcast"
	}))

	feeds := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "feed"
	}))

	roadmaps := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return item.Type == "roadmaps"
	}))

	const readmeTemplate = `
# Awesome Tech Lead [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)

Uma lista de conteúdo sobre lideraça técnica curada pelos membros da comunidade [TechLeads.club](https://comece.techleads.club?utm_source=awesome-tech-lead).

## Pilares

- Excelência Técnica
- Entrega de Valor
- Liderança e Inspiração

{{if .Books}}
## 📚 Livros 

| Título                                                          | Tags  | Nível | Pago? | 
|-----------------------------------------------------------------|-------|-------|-------|
{{- range .Books }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Articles}}
## 📰 Artigos

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Articles }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Courses}}
## 🎓 Cursos

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Courses }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Videos}}
## 🎥 Vídeos

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Videos }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Podcasts}}
## 🎙️ Podcasts

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Podcasts }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Feeds}}
## 📡 Feeds

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Feeds }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}

{{if .Roadmaps}}
## 🗺️ Roadmaps

| Título                                                                    | Tags  | Nível | Pago? | 
|---------------------------------------------------------------------------|-------|--------|-------|
{{- range .Roadmaps }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
{{end}}
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer

	templateData := map[string][]FormattedItem{
		"Books":    books,
		"Articles": articles,
		"Courses":  courses,
		"Videos":   videos,
		"Podcasts": podcasts,
		"Feeds":    feeds,
		"Roadmaps": roadmaps,
	}

	if err := tmpl.Execute(&buf, templateData); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func formatCatalogItems(items []catalog.CatalogItem) []FormattedItem {
	var formattedItems []FormattedItem

	for _, item := range items {
		formattedItems = append(formattedItems, FormattedItem{
			Title:  getTitle(item),
			Author: item.Author,
			Type:   translate(item.Type),
			Tags:   safeJoin(item.Tags, ", "),
			Level:  item.Level,
			IsPaid: getFreeBadge(item.IsPaid),
			URL:    item.URL,
		})
	}

	return formattedItems
}

func getTitle(item catalog.CatalogItem) string {
	// Prevent the pipe from breaking the markdown format.
	return strings.ReplaceAll(item.Title, "|", "-")
}

func safeJoin(slice []string, sep string) string {
	if slice == nil {
		return ""
	}

	return strings.Join(slice, sep)
}

func getFreeBadge(isPaid bool) string {
	if isPaid {
		return "❌"
	}

	return "✅"
}
