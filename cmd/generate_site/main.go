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

	fmt.Println("write public/index.html file")
	err = os.WriteFile("public/index.html", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("error writing public/index.html", err)
		os.Exit(1)
	}
}
