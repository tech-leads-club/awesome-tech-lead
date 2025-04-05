.PHONY: generate-readme setup

setup:
	@echo "Instalando dependências..."
	go mod download
	@echo "Dependências instaladas com sucesso!"

generate-readme:
	@go run cmd/generate_readme/main.go

generate-site:
	@go run cmd/generate_site/main.go

serve-site:
	@go run cmd/serve_site/main.go

all: setup generate-readme generate-site
