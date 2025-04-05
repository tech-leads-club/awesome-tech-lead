package catalog

import (
	"fmt"
	"os"
	"sort"
	"text/template"
)

type PageData struct {
	Items   []CatalogItem
	Filters Filters
}

type Filters struct {
	Tags      []string
	Types     []string
	Levels    []string
	Languages []string
	Prices    []string
}

func SiteTmpl() *template.Template {
	funcMap := template.FuncMap{
		"formatLanguage": FormatLanguage,
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		fmt.Println("error parsing template:", err)
		os.Exit(1)
	}

	return tmpl
}

func BuildPageData(items []CatalogItem) PageData {
	filters := extractFilters(items)
	return PageData{
		Items:   items,
		Filters: filters,
	}
}

func extractFilters(items []CatalogItem) Filters {
	filters := Filters{}
	tagMap := make(map[string]bool)
	typeMap := make(map[string]bool)
	levelMap := make(map[string]bool)
	langMap := make(map[string]bool)
	priceMap := make(map[string]bool)

	for _, item := range items {
		for _, tag := range item.Tags {
			tagMap[tag] = true
		}
		typeMap[item.Type] = true
		levelMap[item.Level] = true
		langMap[item.Language] = true

		if item.IsPaid {
			priceMap["Pago"] = true
		} else {
			priceMap["Gratuito"] = true
		}
	}

	filters.Tags = keyMapToSortedSlice(tagMap)
	filters.Types = keyMapToSortedSlice(typeMap)
	filters.Levels = keyMapToSortedSlice(levelMap)
	filters.Languages = keyMapToSortedSlice(langMap)
	filters.Prices = keyMapToSortedSlice(priceMap)

	return filters
}

func keyMapToSortedSlice(m map[string]bool) []string {
	slice := make([]string, 0, len(m))
	for k := range m {
		slice = append(slice, k)
	}
	sort.Strings(slice)
	return slice
}

func FormatLanguage(lang string) string {
	switch lang {
	case "pt_br":
		return "🇧🇷 Português"
	case "en_us":
		return "🇺🇸 Inglês"
	default:
		return lang
	}
}
