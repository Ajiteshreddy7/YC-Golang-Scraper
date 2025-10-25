# YC-Go-JobScraper



> 🚀 An intelligent job aggregator that automatically discovers and tracks early-career opportunities from Y Combinator companies. Built with Go, powered by GitHub Actions, and hosted for free on GitHub Pages.An intelligent job scraper written in Go that discovers early‑career roles from Greenhouse boards, stores them in SQLite, exports CSV, and serves a static dashboard via GitHub Pages.



[![Deploy Status](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/workflows/Deploy%20to%20GitHub%20Pages/badge.svg)](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/actions)


[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)- Smart scraping via Greenhouse API

- Early‑career filtering (intern, new grad, junior) and US‑location bias

## 📋 Table of Contents- SQLite storage with de‑duplication on URL

- CSV export to `data/job_applications.csv`

- [Overview](#overview)- Static website generation for GitHub Pages

- [Features](#features)- Automated daily scraping via GitHub Actions

- [Live Demo](#live-demo)- 100% free hosting 

- [Quick Start](#quick-start)

- [Installation](#installation)

- [Configuration](#configuration)

- [Usage](#usage)Visit the live dashboard at: **https://ajiteshreddy7.github.io/YC-Go-Scraper/**

- [Architecture](#architecture)

- [GitHub Pages Deployment](#github-pages-deployment) Updates automatically every day at 3 AM UTC.

# YC Job Scraper

Lightweight job aggregator for YC companies. Scrapes early‑career roles (internships, new grad, junior) from Greenhouse, stores them in SQLite, and publishes a searchable static site on GitHub Pages.

• Live site: https://ajiteshreddy7.github.io/YC-Golang-Scraper/

## What it does

- Scrape YC company jobs via Greenhouse
- Auto‑tag levels (Intern, New Grad, Entry)
- Save to SQLite and export CSV
- Generate a static dashboard with search, filters, and “Mark Applied”

## 🏗️ Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      GitHub Actions                          │
│  (Runs daily at 3 AM UTC or on-demand via manual trigger)  │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
         ┌──────────────────────┐
         │   Scraper Binary     │
         │  (cmd/scraper)       │
         │                      │
         │ • Queries Greenhouse │
         │ • Filters jobs       │
         │ • Deduplicates       │
         └──────────┬───────────┘
                    │
                    ▼
         ┌──────────────────────┐
         │   SQLite Database    │
         │   (data/jobs.db)     │
         │                      │
         │ • Job title          │
         │ • Company            │
         │ • Location           │
         │ • Apply URL          │
         │ • Posted date        │
         └──────────┬───────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
        ▼                       ▼
┌─────────────────┐   ┌─────────────────┐
│ Static Site Gen │   │  CSV Exporter   │
│ (cmd/static)    │   │ (internal/exp.) │
│                 │   │                 │
│ • index.html    │   │ • job_apps.csv  │
│ • jobs.json     │   └─────────────────┘
│ • Search UI     │
│ • Filters       │
└────────┬────────┘
         │
         ▼
┌─────────────────────┐
│   GitHub Pages      │
│                     │
│ • Hosts static site │
│ • Global CDN        │
│ • Free HTTPS        │
└─────────────────────┘

```

### Project Structure

```
YC-Golang-Scraper/
├── .github/
│   └── workflows/
│       └── deploy-pages.yml      # GitHub Actions workflow for automation
│
├── config/
│   └── scraper_config.json       # List of companies to scrape
│
├── data/                         # Generated data (not in git)
│   ├── jobs.db                   # SQLite database
│   └── job_applications.csv      # CSV export
│
├── go-scraper/                   # Go application root
│   ├── cmd/
│   │   ├── scraper/
│   │   │   └── main.go           # CLI scraper entrypoint
│   │   ├── static-site/
│   │   │   └── main.go           # Static site generator
│   │   └── dashboard/
│   │       └── main.go           # Optional local API server
│   │
│   ├── internal/
│   │   ├── db/
│   │   │   └── db.go             # SQLite connection and schema
│   │   ├── scraper/
│   │   │   ├── greenhouse.go     # Greenhouse API client
│   │   │   ├── greenhouse_test.go
│   │   │   └── greenhouse_filters_test.go
│   │   ├── exporter/
│   │   │   └── exporter.go       # CSV export logic
│   │   └── logger/
│   │       └── logger.go         # Structured logging
│   │
│   ├── go.mod                    # Go module definition
│   └── go.sum                    # Dependency checksums
│
├── public/                       # Generated static site (not in git)
│   ├── index.html                # Dashboard UI
│   └── jobs.json                 # JSON API
│
├── .gitignore
└── README.md
```
## Docs

- docs/quickstart.md – local setup (Windows/macOS/Linux)
- docs/deployment.md – GitHub Pages + workflow
- docs/configuration.md – config file and env vars


## 🐛 Troubleshooting

### Common Issues

**Issue: Scraper finds no jobs**
- Check company names in `config/scraper_config.json`
- Verify companies use Greenhouse (visit `boards.greenhouse.io/<company>`)
- Run with `LOG_LEVEL=DEBUG` for detailed output

**Issue: GitHub Pages not deploying**
- Ensure GitHub Pages is enabled (Settings → Pages → Source: GitHub Actions)
- Check Actions tab for deployment logs
- Verify workflow file exists at `.github/workflows/deploy-pages.yml`

**Issue: SQLite database locked**
- Close any open database connections
- Delete `jobs.db-shm` and `jobs.db-wal` files
- Restart the scraper

## 📄 License

This project is licensed under the **MIT License**.


<div align="center">

**[⬆ Back to Top](#yc-job-scraper)**



</div>
