package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	catalog "github.com/tech-leads-club/awesome-tech-lead/internal"
)

const tailwindInstallDoc = "https://tailwindcss.com/blog/standalone-cli"

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

	fmt.Println("[site] compiling css...")

	if _, err := os.Stat("./tailwindcss"); os.IsNotExist(err) {
		fmt.Printf(`
ERROR: Tailwind CSS binary not found.
To install the standalone Tailwind CSS CLI:
1. Download the appropriate binary for your system (%s)
2. Place it in your project root as 'tailwindcss'
3. Make it executable (chmod +x tailwindcss on Unix systems)

Installation documentation: %s

`, runtime.GOOS, tailwindInstallDoc)
		os.Exit(1)
	}

	cmd := exec.Command(
		"./tailwindcss",
		"-i", "public/css/main.css",
		"-o", "build/site/css/main.css",
		"--minify",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("error compiling css: %v\noutput: %s\n", err, output)
		os.Exit(1)
	}

	fmt.Println("[site] css built successfully")
}
