package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

// Job represents a single job application record, shared by the frontend and database.
type Job struct {
	ID        int
	Title     string
	Company   string
	Location  string
	Type      string
	URL       string
	DateAdded time.Time
	Status    string
}

// JobFilter is used to define search and pagination parameters for job listing.
// (Currently a placeholder, but essential for main.go compilation).
type JobFilter struct {
	Status string
	Search string
}

type DB struct {
	Conn *sql.DB
}

// Connect opens a SQLite connection using DB_PATH env var or default.
// Also ensures that required directories exist.
func Connect() (*DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/jobs.db"
	}

	// Create data directory if it doesn't exist
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// SQLite tuning for concurrent read-safe performance
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxLifetime(time.Minute * 5)

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	db := &DB{Conn: conn}

	// Initialize schemas if they don't exist
	if err := db.CreateSchema(); err != nil {
		return nil, err
	}
	if err := db.CreateUserSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

// Close closes the database connection.
func (d *DB) Close() error {
	return d.Conn.Close()
}

// -------------------- JOBS TABLE --------------------

// CreateSchema ensures the job_applications table exists.
func (d *DB) CreateSchema() error {
	q := `
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
	`
	_, err := d.Conn.Exec(q)
	return err
}

// InsertJob inserts a job record using a map, ignores duplicate URLs.
func (d *DB) InsertJob(job map[string]interface{}) error {
	q := `INSERT INTO job_applications(title, company, location, type, url)
			 VALUES($1,$2,$3,$4,$5)
			 ON CONFLICT (url) DO NOTHING;`
	_, err := d.Conn.Exec(q,
		job["Title"], job["Company"], job["Location"], job["Type"], job["URL"])
	return err
}

// InsertJobTyped inserts a job record using typed parameters.
func (d *DB) InsertJobTyped(title, company, location, typ, url string) error {
	q := `INSERT INTO job_applications(title, company, location, type, url)
			 VALUES($1,$2,$3,$4,$5)
			 ON CONFLICT (url) DO NOTHING;`
	_, err := d.Conn.Exec(q, title, company, location, typ, url)
	return err
}

// ListJobs retrieves job records based on filters, for the dashboard display.
func (d *DB) ListJobs(filter JobFilter, page, pageSize int) ([]Job, error) {
	q := `
	SELECT id, title, company, location, type, url, date_added, status
	FROM job_applications
	ORDER BY date_added DESC
	LIMIT $1 OFFSET $2
	`
	offset := (page - 1) * pageSize
	rows, err := d.Conn.Query(q, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		// Note: The Scan order MUST match the SELECT column order
		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Company,
			&job.Location,
			&job.Type,
			&job.URL,
			&job.DateAdded,
			&job.Status,
		)
		if err != nil {
			// Log error and continue if a single row is problematic
			continue
		}
		jobs = append(jobs, job)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return jobs, nil
}

// -------------------- USERS TABLE --------------------

// CreateUserSchema ensures a users table exists for authentication.
func (d *DB) CreateUserSchema() error {
	q := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := d.Conn.Exec(q)
	return err
}

// CreateUser registers a new user with bcrypt password hashing.
func (d *DB) CreateUser(username, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	q := `INSERT INTO users(username, password_hash) VALUES($1,$2)`
	_, err = d.Conn.Exec(q, username, string(hashed))
	return err
}

// AuthenticateUser checks if username/password is valid.
func (d *DB) AuthenticateUser(username, password string) (bool, error) {
	var storedHash string
	err := d.Conn.QueryRow(`SELECT password_hash FROM users WHERE username=$1`, username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // user not found
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		return false, nil // password mismatch
	}
	return true, nil
}

// GetUserByUsername retrieves basic user info by username.
func (d *DB) GetUserByUsername(username string) (int, string, time.Time, error) {
	var (
		id        int
		u         string
		createdAt time.Time
	)
	err := d.Conn.QueryRow(
		`SELECT id, username, created_at FROM users WHERE username=$1`,
		username,
	).Scan(&id, &u, &createdAt)
	return id, u, createdAt, err
}