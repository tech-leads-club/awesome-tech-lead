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

func GenerateREADME(items []catalog.CatalogItem) (string, error) {
	formattedItems := formatCatalogItems(items)

	const readmeTemplate = `
# Awesome Tech Lead [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)

Uma lista de conteúdo sobre lideraça técnica curada pelos membros da comunidade [TechLeads.club](https://comece.techleads.club?utm_source=awesome-tech-lead).

## Pilares

- Excelência Técnica
- Entrega de Valor
- Liderança e Inspiração

## Conteúdo 

| Título                      | Tipo | Tags  | Nível | Pago? | 
|-----------------------------|------|-------|-------|-------|
{{- range . }}
| [{{ .Title }}]({{ .URL }}) | {{ .Type }} | {{ .Tags }} | {{ .Level }} | {{ .IsPaid }} | 
{{- end }}
`

	tmpl, err := template.New("readme").Parse(readmeTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, formattedItems); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func formatCatalogItems(items []catalog.CatalogItem) []FormattedItem {
	var formattedItems []FormattedItem

	for _, item := range items {
		formattedItems = append(formattedItems, FormattedItem{
			Title:  getTitle(item),
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
	if item.Author != nil {
		return fmt.Sprintf("%s (por %s)", item.Title, *item.Author)
	}

	return item.Title
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
