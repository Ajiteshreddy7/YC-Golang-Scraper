# Quick Start

This guide shows how to run YC Job Scraper locally on Windows, macOS, and Linux.

## Prerequisites
- Go 1.21+
- Git

## Windows (PowerShell)

```powershell
# Clone repository
git clone https://github.com/Ajiteshreddy7/YC-Golang-Scraper.git
cd YC-Golang-Scraper\go-scraper

# Dependencies
go mod download

# Run scraper
$env:DB_PATH='../data/jobs.db'; $env:LOG_LEVEL='INFO'; go run ./cmd/scraper --config ../config/scraper_config.json --out ../data/job_applications.csv

# Generate static site
$env:DB_PATH='../data/jobs.db'; go run ./cmd/static-site --out ../public

# Open in browser
cd ..; start public/index.html
```

## macOS / Linux (Bash)

```bash
# Clone repository
git clone https://github.com/Ajiteshreddy7/YC-Golang-Scraper.git
cd YC-Golang-Scraper/go-scraper

# Dependencies
go mod download

# Run scraper
DB_PATH=../data/jobs.db LOG_LEVEL=INFO go run ./cmd/scraper --config ../config/scraper_config.json --out ../data/job_applications.csv

# Generate static site
DB_PATH=../data/jobs.db go run ./cmd/static-site --out ../public

# Open in browser
open ../public/index.html  # macOS
xdg-open ../public/index.html  # Linux
```

## Binaries (optional)

```powershell
# From go-scraper
go build -o scraper.exe ./cmd/scraper
go build -o static-site.exe ./cmd/static-site

$env:DB_PATH='../data/jobs.db'; ./scraper.exe --config ../config/scraper_config.json
$env:DB_PATH='../data/jobs.db'; ./static-site.exe --out ../public
```
