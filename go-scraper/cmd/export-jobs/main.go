package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
)

type JobExport struct {
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	DateAdded time.Time `json:"date_added"`
	Status    string    `json:"status"`
}

func main() {
	// Connect to local database
	database, err := db.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Get all jobs
	jobs, err := database.ListJobs(db.JobFilter{}, 1, 1000)
	if err != nil {
		fmt.Printf("Failed to fetch jobs: %v\n", err)
		os.Exit(1)
	}

	// Convert to export format
	var exportJobs []JobExport
	for _, job := range jobs {
		exportJobs = append(exportJobs, JobExport{
			Title:     job.Title,
			Company:   job.Company,
			Location:  job.Location,
			Type:      job.Type,
			URL:       job.URL,
			DateAdded: job.DateAdded,
			Status:    job.Status,
		})
	}

	// Export to JSON
	jsonData, err := json.MarshalIndent(exportJobs, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	outputFile := "jobs_export.json"
	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Successfully exported %d jobs to %s\n", len(exportJobs), outputFile)
	fmt.Printf("📁 File size: %.2f KB\n", float64(len(jsonData))/1024)
	fmt.Printf("\n🚀 Next steps:")
	fmt.Printf("\n1. Upload this file to a public URL (GitHub gist, pastebin, etc.)")
	fmt.Printf("\n2. Go to: https://yc-golang-scraper.onrender.com/import-jobs?url=YOUR_JSON_URL")
	fmt.Printf("\n3. Your Render deployment will import all jobs!")
}
