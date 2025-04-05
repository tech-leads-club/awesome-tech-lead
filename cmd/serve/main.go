package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

func StartServer(items []catalog.CatalogItem) error {
	funcMap := template.FuncMap{
		"formatLanguage": catalog.FormatLanguage,
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).ParseFiles("templates/index.html")
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := catalog.BuildPageData(items)

		tmpl.Execute(w, data)
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("."))))
	log.Printf("Server starting at http://localhost:8080")
	return http.ListenAndServe(":8080", nil)
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

	err = StartServer(items)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
