# YC-Go-JobScraper



> 🚀 An intelligent job aggregator that automatically discovers and tracks early-career opportunities from Y Combinator companies. Built with Go, powered by GitHub Actions, and hosted for free on GitHub Pages.An intelligent job scraper written in Go that discovers early‑career roles from Greenhouse boards, stores them in SQLite, exports CSV, and serves a static dashboard via GitHub Pages.



[![Deploy Status](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/workflows/Deploy%20to%20GitHub%20Pages/badge.svg)](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/actions)## Features

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)- Smart scraping via Greenhouse API

- Early‑career filtering (intern, new grad, junior) and US‑location bias

## 📋 Table of Contents- SQLite storage with de‑duplication on URL

- CSV export to `data/job_applications.csv`

- [Overview](#overview)- Static website generation for GitHub Pages

- [Features](#features)- Automated daily scraping via GitHub Actions

- [Live Demo](#live-demo)- 100% free hosting (no credit card required)

- [Quick Start](#quick-start)

- [Installation](#installation)## Live Dashboard

- [Configuration](#configuration)

- [Usage](#usage)Visit the live dashboard at: **https://ajiteshreddy7.github.io/YC-Go-Scraper/**

- [Architecture](#architecture)

- [GitHub Pages Deployment](#github-pages-deployment)Updates automatically every day at 3 AM UTC.

- [API Reference](#api-reference)

- [Development](#development)## Quick start (Local)

- [Contributing](#contributing)

- [Troubleshooting](#troubleshooting)```powershell

- [License](#license)# From repo root

cd go-scraper

## 🎯 Overview

# Run scraper

YC Job Scraper is a fully automated job discovery platform that:go run ./cmd/scraper --config ../config/scraper_config.json



- **Scrapes** job postings from Y Combinator portfolio companies using Greenhouse APIs# Generate static site

- **Filters** for early-career positions (internships, new grad, entry-level, junior roles)go run ./cmd/static-site --out ../public

- **Stores** data in a lightweight SQLite database with automatic deduplication

- **Exports** to CSV for easy importing into spreadsheet applications# View locally (open public/index.html in browser)

- **Generates** a beautiful static website with search, filtering, and "Mark Applied" trackingstart ../public/index.html

- **Deploys** automatically to GitHub Pages with daily updates```

- **Costs** $0 - completely free hosting with no credit card required

The SQLite database is stored at `data/jobs.db` and CSV at `data/job_applications.csv`.

Perfect for new grads, career switchers, and anyone seeking opportunities at high-growth YC startups.

## Configuration

## ✨ Features

Edit `config/scraper_config.json`:

### 🔍 Smart Job Discovery

- **Greenhouse API Integration**: Directly queries company Greenhouse boards for real-time data```json

- **Intelligent Filtering**: Automatically detects early-career roles via regex patterns{

- **Location Prioritization**: Surfaces US-based positions while including remote opportunities  "target_platforms": {

- **Auto-deduplication**: Prevents duplicate listings by tracking unique job URLs    "greenhouse": ["stripe", "airbnb", "coinbase", "databricks"]

  }

### 💾 Data Management}

- **SQLite Database**: Fast, reliable, file-based storage requiring no server setup```

- **CSV Export**: One-click export to `data/job_applications.csv` for tracking in Excel/Sheets

- **Persistent Storage**: All job data preserved across scraping runsEnvironment variables:

- **Schema Migrations**: Database schema managed via code for easy updates

- `DB_PATH` (optional) override for SQLite database path; defaults to `data/jobs.db`.

### 🎨 Interactive Dashboard- `LOG_LEVEL` one of `DEBUG, INFO, WARN, ERROR` (default `INFO`).

- **Search & Filter**: Find jobs by keyword, level, company, location, and status

- **Mark Applied**: Track your applications with persistent localStorage## GitHub Pages Deployment

- **Live Statistics**: See matching jobs, applied count, and remaining opportunities

- **Responsive Design**: Works beautifully on desktop, tablet, and mobileThe repository is configured to automatically:

- **CSV Export**: Download filtered results directly from the browser1. Run the scraper daily at 3 AM UTC

- **No Backend Required**: Fully client-side rendering for instant load times2. Generate a static website

3. Deploy to GitHub Pages

### ⚙️ Automation & Deployment

- **GitHub Actions**: Automated daily scraping at 3 AM UTC### Enable GitHub Pages:

- **Zero Configuration**: Works out-of-the-box with sensible defaults1. Go to your repo Settings → Pages

- **GitHub Pages**: Free static hosting with global CDN2. Source: "GitHub Actions"

- **Manual Triggers**: Run scraper on-demand via GitHub Actions UI3. The site will be available at: `https://yourusername.github.io/YC-Go-Scraper/`

- **Workflow Caching**: Speeds up builds by caching Go dependencies

### Manual trigger:

## 🌐 Live Demo- Go to Actions tab → "Deploy to GitHub Pages" → "Run workflow"



**Visit:** [https://ajiteshreddy7.github.io/YC-Golang-Scraper/](https://ajiteshreddy7.github.io/YC-Golang-Scraper/)## JSON API



The dashboard updates automatically every day at **3:00 AM UTC** with fresh job postings.The static site also generates a `jobs.json` file available at:

- `https://ajiteshreddy7.github.io/YC-Go-Scraper/jobs.json`

## 🚀 Quick Start

This contains all job data in JSON format for programmatic access.

### Prerequisites

- **Go 1.21+** installed ([Download](https://go.dev/dl/))## Local development

- **Git** installed ([Download](https://git-scm.com/downloads))

- A **GitHub account** (for deployment)```powershell

# Start from go-scraper directory

### Clone and Run Locallycd go-scraper



#### On Windows (PowerShell):# Build all commands

go build ./cmd/scraper

```powershellgo build ./cmd/static-site

# Clone repositorygo build ./cmd/dashboard  # Optional: for local API server

git clone https://github.com/Ajiteshreddy7/YC-Golang-Scraper.git

cd YC-Golang-Scraper# Run scraper

$env:DB_PATH="data/jobs.db"; ./scraper --config ../config/scraper_config.json

# Navigate to Go workspace

cd go-scraper# Generate static site

$env:DB_PATH="data/jobs.db"; ./static-site --out ../public

# Download dependencies

go mod download# Optional: Run dashboard server locally

$env:DB_PATH="data/jobs.db"; ./dashboard --port 8080

# Run scraper to fetch jobs

$env:DB_PATH='../data/jobs.db'; $env:LOG_LEVEL='INFO'; go run ./cmd/scraper --config ../config/scraper_config.json --out ../data/job_applications.csv# Tests

go test ./...

# Generate static site```

$env:DB_PATH='../data/jobs.db'; go run ./cmd/static-site --out ../public

## Project structure

# Open in browser

start ../public/index.html```

```YC-Go-Scraper/

├─ go-scraper/

#### On macOS/Linux (Bash):│  ├─ cmd/

│  │  ├─ scraper/      # CLI scraper entrypoint

```bash│  │  ├─ static-site/  # Static site generator for GitHub Pages

# Clone repository│  │  └─ dashboard/    # Optional: local API server

git clone https://github.com/Ajiteshreddy7/YC-Golang-Scraper.git│  ├─ internal/

cd YC-Golang-Scraper│  │  ├─ db/           # SQLite layer

│  │  ├─ scraper/      # Greenhouse scraper and filters

# Navigate to Go workspace│  │  ├─ exporter/     # CSV exporter

cd go-scraper│  │  └─ logger/       # Structured logging

│  ├─ go.mod

# Download dependencies│  └─ go.sum

go mod download├─ config/

│  └─ scraper_config.json

# Run scraper to fetch jobs├─ data/                # SQLite DB and CSV output

DB_PATH=../data/jobs.db LOG_LEVEL=INFO go run ./cmd/scraper --config ../config/scraper_config.json --out ../data/job_applications.csv├─ public/              # Generated static site

└─ .github/workflows/

# Generate static site   └─ deploy-pages.yml  # Automated scraping and deployment

DB_PATH=../data/jobs.db go run ./cmd/static-site --out ../public```



# Open in browser## License

open ../public/index.html  # macOS

xdg-open ../public/index.html  # LinuxMIT
```

That's it! You now have a local copy with job data.

## 📦 Installation

### Method 1: Build Executables

```powershell
cd go-scraper

# Build all binaries
go build -o scraper.exe ./cmd/scraper
go build -o static-site.exe ./cmd/static-site
go build -o dashboard.exe ./cmd/dashboard

# Run built executables
$env:DB_PATH='../data/jobs.db'; ./scraper.exe --config ../config/scraper_config.json
$env:DB_PATH='../data/jobs.db'; ./static-site.exe --out ../public
```

### Method 2: Use `go run` (No Build)

```powershell
cd go-scraper

# Run directly without building
go run ./cmd/scraper --config ../config/scraper_config.json
go run ./cmd/static-site --out ../public
```

### Method 3: Install Globally

```powershell
cd go-scraper

# Install to $GOPATH/bin
go install ./cmd/scraper
go install ./cmd/static-site
go install ./cmd/dashboard

# Run from anywhere
scraper --config /path/to/scraper_config.json
static-site --out /path/to/public
```

## ⚙️ Configuration

### Scraper Configuration File

Edit `config/scraper_config.json`:

```json
{
  "target_platforms": {
    "greenhouse": [
      "stripe",
      "airbnb",
      "coinbase",
      "databricks",
      "lattice",
      "retool",
      "anthropic",
      "scale",
      "rippling"
    ]
  }
}
```

**Adding More Companies:**
1. Find the company's Greenhouse board URL (e.g., `boards.greenhouse.io/stripe`)
2. Extract the company slug (e.g., `stripe`)
3. Add to the `greenhouse` array in the config

### Environment Variables

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_PATH` | SQLite database file location | `data/jobs.db` | `../data/jobs.db` |
| `LOG_LEVEL` | Logging verbosity | `INFO` | `DEBUG`, `INFO`, `WARN`, `ERROR` |

**Example Usage:**

```powershell
# Windows PowerShell
$env:DB_PATH='custom/path/jobs.db'; $env:LOG_LEVEL='DEBUG'; ./scraper.exe

# Linux/macOS Bash
DB_PATH=custom/path/jobs.db LOG_LEVEL=DEBUG ./scraper
```

## 🎮 Usage

### Running the Scraper

```powershell
# Basic usage
./scraper --config ../config/scraper_config.json

# With CSV export
./scraper --config ../config/scraper_config.json --out ../data/job_applications.csv

# With custom database path
$env:DB_PATH='../data/custom.db'; ./scraper --config ../config/scraper_config.json

# With debug logging
$env:LOG_LEVEL='DEBUG'; ./scraper --config ../config/scraper_config.json
```

**What the scraper does:**
1. Reads company list from `scraper_config.json`
2. Queries each company's Greenhouse API
3. Filters for early-career positions (intern, new grad, junior, entry-level)
4. Prioritizes US locations but includes remote/international roles
5. Stores jobs in SQLite with deduplication by URL
6. Optionally exports to CSV

### Generating the Static Site

```powershell
# Generate site in public/ directory
./static-site --out ../public

# Generate with custom database
$env:DB_PATH='../data/custom.db'; ./static-site --out ../public
```

**What the site generator does:**
1. Reads all jobs from SQLite database
2. Derives job levels from titles (Intern, New Grad, Entry Level, Junior, Mid-level, Senior)
3. Extracts unique companies and locations for filter dropdowns
4. Generates `index.html` with embedded search/filter UI
5. Generates `jobs.json` with all job data for programmatic access

### Running the Dashboard Server (Optional)

For local development, you can run an optional API server:

```powershell
# Start server on default port 8080
./dashboard

# Custom port
./dashboard --port 3000

# Visit in browser
# http://localhost:8080
```

**API Endpoints:**
- `GET /api/jobs` - Returns all jobs as JSON
- `GET /` - Serves the dashboard UI

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

## 🚀 GitHub Pages Deployment

### Initial Setup

1. **Fork or Clone** this repository to your GitHub account

2. **Create Workflow File** via GitHub UI:
   - Go to your repository on GitHub
   - Click "Add file" → "Create new file"
   - Path: `.github/workflows/deploy-pages.yml`
   - Paste the workflow content (see repository for template)
   - Commit the file

3. **Enable GitHub Pages:**
   - Go to **Settings → Pages**
   - Under **Source**, select **"GitHub Actions"**
   - Click **Save**

4. **Trigger First Deployment:**
   - Go to **Actions** tab
   - Select **"Deploy to GitHub Pages"** workflow
   - Click **"Run workflow"** button
   - Wait 2-3 minutes for completion

5. **Access Your Site:**
   - Visit: `https://yourusername.github.io/YC-Golang-Scraper/`
   - Jobs will update automatically every day at 3 AM UTC

### Manual Deployment

**Trigger on-demand:**
1. Go to **Actions** tab in your repository
2. Click **"Deploy to GitHub Pages"** workflow
3. Click **"Run workflow"** dropdown
4. Select branch (usually `main`)
5. Click green **"Run workflow"** button

**Deployment takes ~2-3 minutes:**
- Checkout code ✓
- Setup Go environment ✓
- Build scraper binary ✓
- Run scraper (fetch jobs) ✓
- Build static-site binary ✓
- Generate HTML/JSON ✓
- Deploy to Pages ✓

## 📡 API Reference

### JSON API

The static site generates a `jobs.json` file for programmatic access.

**Endpoint:** `https://ajiteshreddy7.github.io/YC-Golang-Scraper/jobs.json`

**Response Format:**

```json
[
  {
    "id": 1,
    "title": "Software Engineer, New Grad",
    "company": "Stripe",
    "location": "San Francisco, CA",
    "absolute_url": "https://boards.greenhouse.io/stripe/jobs/123456",
    "updated_at": "2025-10-25T10:30:00Z",
    "level": "New Grad"
  }
]
```

**Job Levels:**
- `Intern` - Internship positions
- `New Grad` - New graduate roles
- `Entry Level` - 0-2 years experience
- `Junior` - 1-3 years experience
- `Mid-level` - 3-5 years experience
- `Senior` - 5+ years experience

### Using the API

**Fetch with JavaScript:**

```javascript
fetch('https://ajiteshreddy7.github.io/YC-Golang-Scraper/jobs.json')
  .then(res => res.json())
  .then(jobs => {
    console.log(`Found ${jobs.length} jobs`);
    const internships = jobs.filter(j => j.level === 'Intern');
    console.log(`${internships.length} internships available`);
  });
```

**Fetch with Python:**

```python
import requests

response = requests.get('https://ajiteshreddy7.github.io/YC-Golang-Scraper/jobs.json')
jobs = response.json()

print(f"Found {len(jobs)} jobs")
new_grad_jobs = [j for j in jobs if j['level'] == 'New Grad']
print(f"{len(new_grad_jobs)} new grad positions")
```

## 🛠️ Development

### Running Tests

```powershell
cd go-scraper

# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

### Adding New Companies

1. Find the company's Greenhouse board (e.g., `boards.greenhouse.io/company-name`)
2. Add the slug to `config/scraper_config.json`
3. Run the scraper to test

### Database Schema

```sql
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    company TEXT,
    location TEXT,
    absolute_url TEXT UNIQUE NOT NULL,
    updated_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## 🤝 Contributing

Contributions are welcome! Here's how you can help:

### Reporting Bugs

Open an issue with:
- Clear description of the bug
- Steps to reproduce
- Expected vs actual behavior
- System info (OS, Go version)

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

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

```
MIT License

Copyright (c) 2025 Ajitesh Reddy Tippireddy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

## 🙏 Acknowledgments

- **Y Combinator** for supporting amazing startups
- **Greenhouse** for providing accessible job APIs
- **GitHub** for free hosting and automation
- **Go Community** for excellent tooling and libraries

---

<div align="center">

**[⬆ Back to Top](#yc-job-scraper)**

Made with ❤️ by [Ajitesh Reddy](https://github.com/Ajiteshreddy7)

</div>
