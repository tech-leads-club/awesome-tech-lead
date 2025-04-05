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

	funcMap := template.FuncMap{
		"formatLanguage": catalog.FormatLanguage,
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("error parsing template:", err)
		os.Exit(1)
	}

	pageData := catalog.BuildPageData(items)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, pageData); err != nil {
		fmt.Println("error executing template:", err)
		os.Exit(1)
	}

	fmt.Println("write public/index.html file")
	err = os.WriteFile("public/index.html", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("error writing public/index.html", err)
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
		return hasTag(item, catalog.TechnicalExcellenceTag)
	}))

	leadershipAndInspiration := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return hasTag(item, catalog.LeadershipAndInspirationTag)
	}))

	deliveringValue := formatCatalogItems(filterItem(items, func(item catalog.CatalogItem) bool {
		return hasTag(item, catalog.DeliveringValueTag)
	}))

	const readmeTemplate = `
# Awesome Tech Lead [![Awesome](https://awesome.re/badge.svg)](https://awesome.re)

Lista de conteúdo sobre lideraça técnica curada pelos membros da comunidade
[TechLeads.club 💎](https://comece.techleads.club?utm_source=awesome-tech-lead&utm_medium=readme).

O conteúdo está categorizado nos três pilares da comunidade: Excelência
Técnica, Entrega de Valor e Liderança e Inspiração.

## 🗂️ Índice

- [🏆 Excelência Técnica](#excelencia-tecnica)
- [📦 Entrega de Valor](#entrega-de-valor)
- [🤝 Liderança e Inspiração](#lideranca-e-inspiracao)
- [🎽 Como Contribuir?](#como-contribuir)

{{if .TechnicalExcellence}}
<h2 id="excelencia-tecnica">🏆 Excelência Técnica</h2>

Pilar focado no domínio e aplicação eficaz de tecnologias, práticas e
arquiteturas para criar soluções robustas, escaláveis e de alta qualidade.

| Título                                                          | Tags        | 
|-----------------------------------------------------------------|-------------|
{{- range .TechnicalExcellence }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}

{{if .DeliveringValue}}
<h2 id="entrega-de-valor">📦 Entrega de Valor</h2>

Pilar relacionado a práticas ágeis e à capacidade de entregar projetos de
software de maneira eficiente, com alinhamento estratégico e foco nas
necessidades do negócio.

| Título                                                          | Tags        |
|-----------------------------------------------------------------|-------------|
{{- range .DeliveringValue }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}

{{if .LeadershipAndInspiration}}
<h2 id="lideranca-e-inspiracao">🤝 Liderança e Inspiração</h2>

Pilar focado na habilidade de liderar times, pessoas, alinhar expectativas, mentorar colegas e
dar feedback.

| Título                                                          | Tags        |
|-----------------------------------------------------------------|-------------|
{{- range .LeadershipAndInspiration }}
| [{{ .Title }}]({{ .URL }}){{if .Author}} por {{.Author}}{{end}} | {{ .Tags }} |
{{- end }}
{{end}}

<h2 id="como-contribuir">🎽 Como Contribuir?</h2>

Deseja contribuir com esse repositório? Saiba mais em
[CONTRIBUTING.md](./CONTRIBUTING.md).
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
