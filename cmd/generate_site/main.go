package main

import (
	"bytes"
	"fmt"
	"os"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

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

	tmpl := catalog.SiteTmpl()
	pageData := catalog.BuildPageData(items)

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, pageData); err != nil {
		fmt.Println("error executing template:", err)
		os.Exit(1)
	}

	fmt.Println("[site] deleting build directory...")
	if err := os.RemoveAll("build"); err != nil {
		fmt.Println("error deleting build directory:", err)
		os.Exit(1)
	}

	fmt.Println("[site] building...")
	if err := os.MkdirAll("build/site", 0755); err != nil {
		fmt.Println("error creating build directory:", err)
		os.Exit(1)
	}

	fmt.Println("[site] copying public/ to build/site")
	if err := os.CopyFS("build/site", os.DirFS("public")); err != nil {
		fmt.Println("error copying public directory:", err)
		os.Exit(1)
	}

	fmt.Println("[site] writing build/site/index.html")
	if err := os.WriteFile("build/site/index.html", buf.Bytes(), 0644); err != nil {
		fmt.Println("error writing build/site/index.html", err)
		os.Exit(1)
	}
}
