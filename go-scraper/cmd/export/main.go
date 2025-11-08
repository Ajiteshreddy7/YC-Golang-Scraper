package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
)

type Job struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	DateAdded time.Time `json:"date_added"`
	Status    string    `json:"status"`
}

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>YC Job Tracker - Static Export</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: Arial, sans-serif; background: #f5f5f5; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .header { text-align: center; margin-bottom: 30px; }
        .header h1 { color: #333; margin-bottom: 10px; }
        .stats { display: flex; justify-content: center; gap: 30px; margin-bottom: 30px; }
        .stat { text-align: center; padding: 15px; background: #f8f9fa; border-radius: 8px; }
        .stat h3 { font-size: 24px; color: #007bff; margin-bottom: 5px; }
        .search-box { margin-bottom: 20px; }
        .search-box input { width: 100%; padding: 12px; border: 2px solid #ddd; border-radius: 8px; font-size: 16px; }
        .job-grid { display: grid; gap: 20px; }
        .job-card { border: 1px solid #ddd; border-radius: 8px; padding: 20px; background: white; transition: transform 0.2s; }
        .job-card:hover { transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.1); }
        .job-title { font-size: 18px; font-weight: bold; color: #333; margin-bottom: 8px; }
        .job-company { font-size: 16px; color: #007bff; margin-bottom: 5px; }
        .job-location { color: #666; margin-bottom: 5px; }
        .job-type { display: inline-block; background: #e9ecef; color: #495057; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-bottom: 10px; }
        .job-date { color: #999; font-size: 14px; margin-bottom: 10px; }
        .job-status { padding: 6px 12px; border-radius: 4px; font-size: 12px; font-weight: bold; }
        .status-applied { background: #d4edda; color: #155724; }
        .status-not-applied { background: #f8d7da; color: #721c24; }
        .job-link { display: inline-block; background: #007bff; color: white; text-decoration: none; padding: 8px 16px; border-radius: 4px; margin-top: 10px; }
        .job-link:hover { background: #0056b3; }
        .note { text-align: center; color: #666; margin-top: 30px; padding: 20px; background: #f8f9fa; border-radius: 8px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 YC Job Tracker</h1>
            <p>Track your Y Combinator job applications</p>
        </div>

        <div class="stats">
            <div class="stat">
                <h3>{{.TotalJobs}}</h3>
                <p>Total Jobs</p>
            </div>
            <div class="stat">
                <h3>{{.Applied}}</h3>
                <p>Applied</p>
            </div>
            <div class="stat">
                <h3>{{.NotApplied}}</h3>
                <p>Not Applied</p>
            </div>
        </div>

        <div class="search-box">
            <input type="text" id="searchInput" placeholder="Search jobs by title, company, or location..." onkeyup="filterJobs()">
        </div>

        <div class="job-grid" id="jobGrid">
            {{range .Jobs}}
            <div class="job-card" data-title="{{.Title}}" data-company="{{.Company}}" data-location="{{.Location}}">
                <div class="job-title">{{.Title}}</div>
                <div class="job-company">{{.Company}}</div>
                <div class="job-location">📍 {{.Location}}</div>
                <div class="job-type">{{.Type}}</div>
                <div class="job-date">📅 {{.DateAdded.Format "Jan 2, 2006"}}</div>
                <div class="job-status {{if eq .Status "Applied"}}status-applied{{else}}status-not-applied{{end}}">
                    {{.Status}}
                </div>
                <a href="{{.URL}}" target="_blank" class="job-link">View Job</a>
            </div>
            {{end}}
        </div>

        <div class="note">
            <p><strong>📝 Note:</strong> This is a static export of your job tracker. To mark jobs as applied or add new jobs, use the local dashboard at <code>localhost:8080</code></p>
            <p>Last updated: {{.LastUpdated}}</p>
        </div>
    </div>

    <script>
        function filterJobs() {
            const searchTerm = document.getElementById('searchInput').value.toLowerCase();
            const jobCards = document.querySelectorAll('.job-card');
            
            jobCards.forEach(card => {
                const title = card.getAttribute('data-title').toLowerCase();
                const company = card.getAttribute('data-company').toLowerCase();
                const location = card.getAttribute('data-location').toLowerCase();
                
                if (title.includes(searchTerm) || company.includes(searchTerm) || location.includes(searchTerm)) {
                    card.style.display = 'block';
                } else {
                    card.style.display = 'none';
                }
            });
        }
    </script>
</body>
</html>`

func main() {
	// Connect to database
	database, err := db.Connect()
	if err != nil {
		fmt.Printf("Database connection failed: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all jobs
	jobs, err := database.ListJobs(db.JobFilter{}, 1, 1000) // Get first 1000 jobs
	if err != nil {
		fmt.Printf("Failed to fetch jobs: %v\n", err)
		os.Exit(1)
	}

	// Convert to template format
	var templateJobs []Job
	applied := 0
	for _, job := range jobs {
		templateJobs = append(templateJobs, Job{
			ID:        job.ID,
			Title:     job.Title,
			Company:   job.Company,
			Location:  job.Location,
			Type:      job.Type,
			URL:       job.URL,
			DateAdded: job.DateAdded,
			Status:    job.Status,
		})
		if job.Status == "Applied" {
			applied++
		}
	}

	// Prepare template data
	data := struct {
		Jobs        []Job
		TotalJobs   int
		Applied     int
		NotApplied  int
		LastUpdated string
	}{
		Jobs:        templateJobs,
		TotalJobs:   len(templateJobs),
		Applied:     applied,
		NotApplied:  len(templateJobs) - applied,
		LastUpdated: time.Now().Format("January 2, 2006 at 3:04 PM"),
	}

	// Parse and execute template
	tmpl, err := template.New("jobs").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("Template parsing failed: %v\n", err)
		os.Exit(1)
	}

	// Create output directory
	outputDir := "../docs"
	os.MkdirAll(outputDir, 0755)

	// Create output file
	outputFile := filepath.Join(outputDir, "index.html")
	file, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Execute template
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("Template execution failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Static site generated successfully!\n")
	fmt.Printf("📄 Output: %s\n", outputFile)
	fmt.Printf("📊 Generated page with %d jobs (%d applied, %d not applied)\n",
		data.TotalJobs, data.Applied, data.NotApplied)
	fmt.Printf("\n🚀 To deploy to GitHub Pages:")
	fmt.Printf("\n   1. Copy the generated file to your repository root")
	fmt.Printf("\n   2. Commit and push to GitHub")
	fmt.Printf("\n   3. Your jobs will be visible at: https://ajiteshreddy7.github.io/YC-Golang-Scraper/\n")
}
