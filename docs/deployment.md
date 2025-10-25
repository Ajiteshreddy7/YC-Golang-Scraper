# Deployment (GitHub Pages)

This project publishes a static site with GitHub Actions + GitHub Pages.

## 1) Create workflow file

Create `.github/workflows/deploy-pages.yml` in your repo with:

```yaml
name: Deploy to GitHub Pages
on:
  schedule:
    - cron: "0 3 * * *"  # daily 03:00 UTC
  workflow_dispatch:
  push:
    branches: [ main ]
permissions:
  contents: read
  pages: write
  id-token: write
jobs:
  scrape-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21.x'
      - name: Download deps
        working-directory: go-scraper
        run: go mod download
      - name: Build scraper
        working-directory: go-scraper
        run: go build ./cmd/scraper
      - name: Build static site
        working-directory: go-scraper
        run: go build ./cmd/static-site
      - name: Run scraper
        working-directory: go-scraper
        env:
          LOG_LEVEL: INFO
          DB_PATH: ../data/jobs.db
        run: ./scraper --config ../config/scraper_config.json --out ../data/job_applications.csv
      - name: Generate site
        working-directory: go-scraper
        env:
          LOG_LEVEL: INFO
          DB_PATH: ../data/jobs.db
        run: ./static-site --out ../public
      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: './public'
      - name: Deploy
        uses: actions/deploy-pages@v4
```

Notes:
- We set `working-directory` per step to `go-scraper`.
- Cache is optional; omitting it avoids noisy warnings.

## 2) Enable Pages
- Settings → Pages → Source: GitHub Actions

## 3) Run first deploy
- Actions tab → Deploy to GitHub Pages → Run workflow

Your site will appear at:
```
https://<your-username>.github.io/YC-Golang-Scraper/
```
