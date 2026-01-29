# YC-Go-JobScraper

## 🚀 Featured Project: YC-Go-JobScraper
[![Showcase](https://img.shields.io/badge/Featured%20on-AI%20Demos-blueviolet)](https://aidemos.com/showcase/447)

> An intelligent job aggregator with authentication that automatically discovers and tracks early-career opportunities from Y Combinator companies. Built with Go, featuring secure user accounts, job tracking, and deployed on Render + GitHub Pages.

[![Deploy Status](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/workflows/Deploy%20to%20GitHub%20Pages/badge.svg)](https://github.com/Ajiteshreddy7/YC-Golang-Scraper/actions)
[![Go Version](https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go)](https://go.dev/)

## 🌟 Features

- **Smart Scraping**: Greenhouse API integration with Y Combinator companies
- **Authentication System**: Secure user accounts with bcrypt password hashing
- **Job Tracking**: Mark jobs as applied, search, filter, and manage applications
- **Dual Deployment**: 
  - 🔐 **Authenticated Dashboard**: [yc-golang-scraper.onrender.com](https://yc-golang-scraper.onrender.com)
  - 📊 **Static Job Listings**: [ajiteshreddy7.github.io/YC-Go-Scraper](https://ajiteshreddy7.github.io/YC-Go-Scraper)
- **Early-Career Focus**: Filters for internships, new grad, and junior positions
- **SQLite Storage**: Local database with deduplication
- **CSV Export**: Download your job applications data
- **Automated Updates**: Daily scraping via GitHub Actions

## 🚀 Live Demo

### 🔐 Authenticated Platform (Full Features)
**URL**: https://yc-golang-scraper.onrender.com

**Features**:
- ✅ Create account & secure login
- ✅ Mark jobs as applied
- ✅ Personal application tracking
- ✅ Search & filter jobs
- ✅ CSV export of applications
- ✅ Responsive design

**Demo Credentials**: `admin` / `password123`

### 📊 Static Job Listings (Public)
**URL**: https://ajiteshreddy7.github.io/YC-Go-Scraper

**Features**:
- ✅ Browse all scraped jobs
- ✅ Search and filter
- ✅ No account required
- ✅ Updated daily at 3 AM UTC



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

### Hybrid Deployment Strategy

```
                    ┌─────────────────────────────────────┐
                    │            GitHub Actions           │
                    │     (Daily scraping at 3 AM UTC)   │
                    └─────────────┬───────────────────────┘
                                  │
                                  ▼
                    ┌─────────────────────────────────────┐
                    │         Scraper Process             │
                    │  • Queries Greenhouse APIs          │
                    │  • Filters Y Combinator companies   │
                    │  • Stores in SQLite database        │
                    └─────────────┬───────────────────────┘
                                  │
                ┌─────────────────┴─────────────────┐
                │                                   │
                ▼                                   ▼
   ┌──────────────────────────┐         ┌──────────────────────────┐
   │     Render.com           │         │     GitHub Pages         │
   │   (Authentication)       │         │    (Static Listings)     │
   │                          │         │                          │
   │ 🔐 User Accounts         │         │ 📊 Public Job Browser    │
   │ 🎯 Job Tracking          │         │ 🔍 Search & Filter       │
   │ 📝 Mark as Applied       │         │ 📱 Mobile Responsive     │
   │ 📊 Personal Dashboard    │         │ 🚀 Global CDN           │
   │ 📥 CSV Export            │         │ ⚡ Fast Loading          │
   │                          │         │                          │
   │ yc-golang-scraper        │         │ ajiteshreddy7.github.io  │
   │ .onrender.com            │         │ /YC-Go-Scraper           │
   └──────────────────────────┘         └──────────────────────────┘
```

### Why Hybrid Architecture?

**🔐 Render (Authentication Platform)**
- Persistent user sessions and secure login
- Personal job application tracking
- Private user data and preferences
- Full CRUD operations on job applications

**📊 GitHub Pages (Public Listings)**
- Fast, global CDN delivery
- No server costs or maintenance
- Great for public job discovery
- SEO-friendly static content

### Project Structure

```
YC-Golang-Scraper/
├── .github/
│   └── workflows/
│       └── deploy-pages.yml      # Auto-deployment workflow
│
├── config/
│   └── scraper_config.json       # Y Combinator companies list
│
├── go-scraper/                   # Go application root
│   ├── cmd/
│   │   ├── dashboard/
│   │   │   └── main.go           # Authentication web server (Render)
│   │   ├── scraper/
│   │   │   └── main.go           # Job scraper CLI
│   │   ├── export/
│   │   │   └── main.go           # Static site generator (GitHub Pages)
│   │   └── export-jobs/
│   │       └── main.go           # Job data export utility
│   │
│   ├── internal/
│   │   ├── db/
│   │   │   └── db.go             # SQLite schema & operations
│   │   ├── scraper/
│   │   │   └── greenhouse.go     # Greenhouse API client
│   │   └── logger/
│   │       └── logger.go         # Structured logging
│   │
│   ├── render.yaml               # Render deployment config
│   ├── go.mod                    # Go module definition
│   └── go.sum                    # Dependencies
│
├── data/
│   ├── jobs.db                   # SQLite database (local)
│   └── job_applications.csv      # CSV exports
│
└── README.md
```
## 🚀 Quick Start

### Option 1: Use the Live Platform (Recommended)
1. Visit **https://yc-golang-scraper.onrender.com**
2. Create an account or use demo login: `admin` / `password123`
3. Start tracking your job applications immediately!

### Option 2: Local Development
```bash
# Clone the repository
git clone https://github.com/Ajiteshreddy7/YC-Golang-Scraper.git
cd YC-Golang-Scraper/go-scraper

# Install Go dependencies
go mod tidy

# Initialize database
go run ./cmd/init

# Run the scraper
go run ./cmd/scraper --config ../config/scraper_config.json

# Start the authentication dashboard
go run ./cmd/dashboard --port 8080
```

Visit http://localhost:8080 to access your local instance.

## 📖 Documentation

- **[Setup Guide](docs/SETUP_GUIDE_V3.md)** - Complete local development setup
- **[Dashboard Setup](docs/DASHBOARD_SETUP.md)** - Authentication system configuration
- **[Scraper Setup](docs/SCRAPER_SETUP.md)** - Job scraping configuration

## 🔧 Configuration

The scraper is configured via `config/scraper_config.json`:

```json
{
  "companies": [
    "stripe",
    "openai", 
    "anthropic",
    "databricks",
    "..."
  ]
}
```

Add Y Combinator portfolio companies to automatically scrape their job boards.

### Automatic import on startup (optional)

The dashboard can automatically import jobs from a public JSON file on startup. Set the environment variable `IMPORT_JOBS_URL` to a publicly-accessible raw JSON URL (for example a GitHub raw URL):

```
IMPORT_JOBS_URL=https://raw.githubusercontent.com/yourname/your-repo/main/jobs_export.json
```

When set, the dashboard will fetch that URL once during startup and import any jobs found. This is handy for one-time bootstrapping of a fresh deployment.

## 🐛 Troubleshooting

### Authentication Issues
**Can't create new accounts**
- Ensure you're using the Render deployment: https://yc-golang-scraper.onrender.com
- Check that registration is enabled (should see "Create account" button)

**Login fails**
- Try demo credentials: `admin` / `password123`
- Clear browser cookies and try again
- Check browser console for JavaScript errors

### Scraper Issues
**No jobs found**
- Verify company names in `config/scraper_config.json` match Greenhouse board names
- Check if companies use Greenhouse: visit `boards.greenhouse.io/<company>`
- Run with `LOG_LEVEL=DEBUG` for detailed logs

**Database locked errors**
- Close any open SQLite connections
- Delete temporary files: `jobs.db-shm`, `jobs.db-wal`
- Restart the application

### Deployment Issues
**GitHub Pages not updating**
- Check Actions tab for workflow status
- Verify Pages settings: Settings → Pages → Source: GitHub Actions
- Workflow file must exist at `.github/workflows/deploy-pages.yml`

**Render deployment failing**
- Check Render dashboard for build logs
- Ensure `render.yaml` configuration is correct
- Verify Go version compatibility (1.21+)


## ⭐ Support

If you found this project helpful, please give it a ⭐ on GitHub!

**🔗 Links:**
- **Authentication Platform**: https://yc-golang-scraper.onrender.com
- **Static Job Listings**: https://ajiteshreddy7.github.io/YC-Go-Scraper
- **GitHub Repository**: https://github.com/Ajiteshreddy7/YC-Golang-Scraper

---

<div align="center">
  <strong>Built with ❤️ for the job search community</strong><br>
  Helping early-career professionals discover opportunities at Y Combinator companies
</div>
