# YC-Go-Scraper

An intelligent job scraper written in Go that discovers early‑career roles from Greenhouse boards, stores them in SQLite, exports CSV, and serves a static dashboard via GitHub Pages.

## Features

- Smart scraping via Greenhouse API
- Early‑career filtering (intern, new grad, junior) and US‑location bias
- SQLite storage with de‑duplication on URL
- CSV export to `data/job_applications.csv`
- Static website generation for GitHub Pages
- Automated daily scraping via GitHub Actions
- 100% free hosting (no credit card required)

## Live Dashboard

Visit the live dashboard at: **https://ajiteshreddy7.github.io/YC-Go-Scraper/**

Updates automatically every day at 3 AM UTC.

## Quick start (Local)

```powershell
# From repo root
cd go-scraper

# Run scraper
go run ./cmd/scraper --config ../config/scraper_config.json

# Generate static site
go run ./cmd/static-site --out ../public

# View locally (open public/index.html in browser)
start ../public/index.html
```

The SQLite database is stored at `data/jobs.db` and CSV at `data/job_applications.csv`.

## Configuration

Edit `config/scraper_config.json`:

```json
{
  "target_platforms": {
    "greenhouse": ["stripe", "airbnb", "coinbase", "databricks"]
  }
}
```

Environment variables:

- `DB_PATH` (optional) override for SQLite database path; defaults to `data/jobs.db`.
- `LOG_LEVEL` one of `DEBUG, INFO, WARN, ERROR` (default `INFO`).

## GitHub Pages Deployment

The repository is configured to automatically:
1. Run the scraper daily at 3 AM UTC
2. Generate a static website
3. Deploy to GitHub Pages

### Enable GitHub Pages:
1. Go to your repo Settings → Pages
2. Source: "GitHub Actions"
3. The site will be available at: `https://yourusername.github.io/YC-Go-Scraper/`

### Manual trigger:
- Go to Actions tab → "Deploy to GitHub Pages" → "Run workflow"

## JSON API

The static site also generates a `jobs.json` file available at:
- `https://ajiteshreddy7.github.io/YC-Go-Scraper/jobs.json`

This contains all job data in JSON format for programmatic access.

## Local development

```powershell
# Start from go-scraper directory
cd go-scraper

# Build all commands
go build ./cmd/scraper
go build ./cmd/static-site
go build ./cmd/dashboard  # Optional: for local API server

# Run scraper
$env:DB_PATH="data/jobs.db"; ./scraper --config ../config/scraper_config.json

# Generate static site
$env:DB_PATH="data/jobs.db"; ./static-site --out ../public

# Optional: Run dashboard server locally
$env:DB_PATH="data/jobs.db"; ./dashboard --port 8080

# Tests
go test ./...
```

## Project structure

```
YC-Go-Scraper/
├─ go-scraper/
│  ├─ cmd/
│  │  ├─ scraper/      # CLI scraper entrypoint
│  │  ├─ static-site/  # Static site generator for GitHub Pages
│  │  └─ dashboard/    # Optional: local API server
│  ├─ internal/
│  │  ├─ db/           # SQLite layer
│  │  ├─ scraper/      # Greenhouse scraper and filters
│  │  ├─ exporter/     # CSV exporter
│  │  └─ logger/       # Structured logging
│  ├─ go.mod
│  └─ go.sum
├─ config/
│  └─ scraper_config.json
├─ data/                # SQLite DB and CSV output
├─ public/              # Generated static site
└─ .github/workflows/
   └─ deploy-pages.yml  # Automated scraping and deployment
```

## License

MIT