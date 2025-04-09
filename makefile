.PHONY: generate-readme setup

setup:
	@echo "Installing golang dependencies"
	go mod download

generate-readme:
	@go run cmd/generate_readme/main.go

# === Site ===

site/generate:
	@go run cmd/generate_site/main.go

SITE_PYTHON := $(shell command -v python3 2> /dev/null || command -v python 2> /dev/null)
SITE_PORT := $(or $(PORT),8080)

site/watch:
	@command -v air >/dev/null 2>&1 || { echo >&2 "Air is not installed. Install with: go install github.com/cosmtrek/air@latest"; exit 1; }
	@echo "Watching changes in the catalog"
	@air

site/serve:
	@make site/watch &

	@if [ -z "$(SITE_PYTHON)" ]; then \
		echo "Error: neither python3 nor python found in PATH"; \
		exit 1; \
	fi
	@echo "Starting server at http://localhost:$(SITE_PORT)"
	@$(SITE_PYTHON) -m http.server $(SITE_PORT) --directory build/site

all: setup generate-readme site/generate
