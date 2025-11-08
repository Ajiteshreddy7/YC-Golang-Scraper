package main

import (
	"database/sql"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func main() {
	// Create data directory
	os.MkdirAll("./data", 0755)

	// Connect to SQLite
	db, err := sql.Open("sqlite", "./data/jobs.db")
	if err != nil {
		fmt.Printf("Failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create tables
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS job_applications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		company TEXT,
		location TEXT,
		salary TEXT,
		type TEXT,
		url TEXT UNIQUE,
		date_added DATETIME DEFAULT CURRENT_TIMESTAMP,
		status TEXT DEFAULT 'Not Applied'
	);
	`)
	if err != nil {
		fmt.Printf("Failed to create job table: %v\n", err)
		os.Exit(1)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		fmt.Printf("Failed to create users table: %v\n", err)
		os.Exit(1)
	}

	// Create admin user
	hashed, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	_, err = db.Exec(`INSERT OR IGNORE INTO users(username, password_hash) VALUES(?, ?)`, "admin", string(hashed))
	if err != nil {
		fmt.Printf("Failed to create admin user: %v\n", err)
		os.Exit(1)
	}

	// Sample job data - insert a few example jobs
	jobs := []struct {
		title, company, location, jobType, url string
	}{
		{"Software Engineer", "Y Combinator", "San Francisco, CA", "Full-time", "https://example.com/job1"},
		{"Full Stack Developer", "OpenAI", "Remote", "Full-time", "https://example.com/job2"},
		{"Backend Engineer", "Stripe", "San Francisco, CA", "Full-time", "https://example.com/job3"},
		{"Frontend Developer", "Airbnb", "San Francisco, CA", "Full-time", "https://example.com/job4"},
		{"DevOps Engineer", "Dropbox", "Remote", "Full-time", "https://example.com/job5"},
	}

	for _, job := range jobs {
		_, err = db.Exec(`INSERT OR IGNORE INTO job_applications(title, company, location, type, url) VALUES(?, ?, ?, ?, ?)`,
			job.title, job.company, job.location, job.jobType, job.url)
		if err != nil {
			fmt.Printf("Failed to insert job %s: %v\n", job.title, err)
		}
	}

	fmt.Printf("✅ Database initialized successfully!\n")
	fmt.Printf("📊 Created admin user and sample jobs\n")
	fmt.Printf("🔐 Username: admin, Password: password123\n")
}
