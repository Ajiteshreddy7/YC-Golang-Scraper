package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/logger"
)

// -------------------- TYPES --------------------

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

// PageData includes 'User' field required for authenticated templates
type PageData struct {
	Jobs       []Job
	TotalJobs  int
	NotApplied int
	Applied    int
	User       string
}

// -------------------- AUTHENTICATION SETUP --------------------

var store = sessions.NewCookieStore([]byte("replace-this-secret-key-123"))

// AuthRequired is middleware to protect authenticated routes
func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if session.Values["user"] == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// -------------------- AUTHENTICATION TEMPLATES --------------------

var loginHTML = `<!DOCTYPE html>
<html>
<head><title>Login</title><style>body {font-family:sans-serif; margin:50px;} h2 {color:#333;}</style></head>
<body>
<h2>Login</h2>
<form method="POST" action="/login">
  <label>Username:</label><br/>
  <input type="text" name="username" required/><br/><br/>
  <label>Password:</label><br/>
  <input type="password" name="password" required/><br/><br/>
  <button type="submit">Login</button>
</form>
<p>No account? <a href="/register">Register here</a></p>
</body>
</html>`

var registerHTML = `<!DOCTYPE html>
<html>
<head><title>Register</title><style>body {font-family:sans-serif; margin:50px;} h2 {color:#333;}</style></head>
<body>
<h2>Register</h2>
<form method="POST" action="/register">
  <label>Username:</label><br/>
  <input type="text" name="username" required/><br/><br/>
  <label>Password:</label><br/>
  <input type="password" name="password" required/><br/><br/>
  <button type="submit">Create Account</button>
</form>
<p>Already have an account? <a href="/login">Login here</a></p>
</body>
</html>`

// -------------------- UI TEMPLATES (from main (1).go) --------------------

const dashboardHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Job Application Dashboard</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            margin: 40px;
            background-color: #f8f9fa;
        }
        h1 {
            color: #343a40;
            text-align: center;
        }
        .stats {
            display: flex;
            justify-content: space-around;
            margin: 30px 0;
        }
        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
            min-width: 150px;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #007bff;
        }
        .stat-label {
            color: #6c757d;
            margin-top: 5px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            background: white;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            border-radius: 8px;
            overflow: hidden;
        }
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #dee2e6;
        }
        th {
            background-color: #007bff;
            color: white;
        }
        tr:hover {
            background-color: #f8f9fa;
        }
        a {
            color: #007bff;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        .status-not-applied {
            color: #dc3545;
            font-weight: bold;
        }
        .status-applied {
            color: #28a745;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom: 20px; border-bottom: 1px solid #ccc; padding-bottom: 10px;">
        <h1 style="margin:0; color: #343a40;">🚀 Job Dashboard</h1>
        <div style="font-size: 0.95em; color: #6c757d;">Logged in as <strong>{{.User}}</strong> — <a href="/filters" style="color: #007bff;">Filters</a> | <a href="/logout" style="color: #dc3545;">Logout</a></div>
    </div>
    
    <div class="stats">
        <div class="stat-card">
            <div class="stat-number">{{.TotalJobs}}</div>
            <div class="stat-label">Total Jobs</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">{{.NotApplied}}</div>
            <div class="stat-label">Not Applied</div>
        </div>
        <div class="stat-card">
            <div class="stat-number">{{.Applied}}</div>
            <div class="stat-label">Applied</div>
        </div>
    </div>
    
    <table>
        <thead>
            <tr>
                <th>Date</th>
                <th>Company</th>
                <th>Title</th>
                <th>Location</th>
                <th>Status</th>
                <th>Link</th>
            </tr>
        </thead>
        <tbody>
            {{range .Jobs}}
            <tr>
                <td>{{.DateAdded.Format "2006-01-02"}}</td>
                <td>{{.Company}}</td>
                <td>{{.Title}}</td>
                <td>{{.Location}}</td>
                <td><span class="status-{{.Status | lower}}">{{.Status}}</span></td>
                <td><a href="{{.URL}}" target="_blank">Apply</a></td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>
`

const landingHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Find Jobs</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 40px; background: #f8f9fa; }
		.container { max-width: 900px; margin: 0 auto; }
		h1 { color: #343a40; text-align: center; }
		form { background: #fff; padding: 24px; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,.08); }
		.section { margin-bottom: 20px; }
		.section-title { font-weight: 600; color: #495057; margin-bottom: 8px; }
		.levels { display: grid; grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); gap: 8px 16px; margin: 12px 0 20px; }
		.level { background:#f1f3f5; padding:10px 12px; border-radius: 6px; }
		.actions { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; }
		input[type="text"] { padding: 10px 12px; border:1px solid #ced4da; border-radius: 6px; width: 260px; }
		select { padding: 10px 12px; border:1px solid #ced4da; border-radius:6px; min-width: 200px; }
		button { background:#007bff; color:#fff; border:none; padding:10px 16px; border-radius:6px; cursor:pointer; font-weight: 500; }
		button:hover { background:#0069d9; }
		.note { color:#6c757d; font-size: 0.95em; margin-top: 8px; }
	</style>
	<script>
	  function toggleAll(source) {
		const checkboxes = document.querySelectorAll('input[name="level"]');
		checkboxes.forEach(cb => cb.checked = source.checked);
	  }
	</script>
 </head>
 <body>
   <div style="max-width: 900px; margin: 0 auto; display:flex; justify-content:flex-end; padding-bottom: 10px;">
        <div style="font-size: 0.95em; color: #6c757d;">Logged in as <strong>{{.User}}</strong> — <a href="/dashboard">Dashboard</a> | <a href="/logout">Logout</a></div>
    </div>
   <div class="container">
	 <h1>What jobs are you looking for?</h1>
	 <form method="GET" action="/results">
		<div class="section">
		  <div class="section-title">Job Levels</div>
		  <label><input type="checkbox" onclick="toggleAll(this)"> Select/Deselect All</label>
		  <div class="levels">
			 {{range .Levels}}
			 <label class="level"><input type="checkbox" name="level" value="{{.}}"> {{.}}</label>
			 {{end}}
		  </div>
		</div>
		<div class="section">
		  <div class="section-title">Additional Filters</div>
		  <div class="actions">
			 <input type="text" name="q" placeholder="Search titles..." />
			 <select name="company">
				<option value="">All Companies</option>
				{{range .Companies}}
				<option value="{{.}}">{{.}}</option>
				{{end}}
			 </select>
			 <select name="location">
				<option value="">All Locations</option>
				{{range .Locations}}
				<option value="{{.}}">{{.}}</option>
				{{end}}
			 </select>
			 <select name="status">
				<option value="">All Statuses</option>
				<option value="Not Applied">Not Applied</option>
				<option value="Applied">Applied</option>
			 </select>
			 <button type="submit">Show Jobs</button>
		  </div>
		</div>
		<div class="note">Levels are detected from existing job titles in your database.</div>
	 </form>
   </div>
 </body>
 </html>
`

const resultsHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Results</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 40px; background: #f8f9fa; }
		.container { max-width: 1100px; margin: 0 auto; }
		h1 { color: #343a40; }
		.pill { display:inline-block; background:#e9ecef; color:#495057; padding:6px 10px; border-radius:999px; margin: 0 6px 6px 0; font-size: 0.9em; }
		ul { list-style: none; padding: 0; }
		li { background:#fff; padding:14px 16px; border-radius:8px; margin-bottom:10px; box-shadow: 0 2px 6px rgba(0,0,0,.06); display:flex; justify-content:space-between; align-items:center; gap: 10px; }
		.meta { color:#6c757d; font-size: 0.95em; }
		a.btn { background:#007bff; color:#fff; padding:8px 12px; border-radius:6px; text-decoration:none; display:inline-block; margin-right: 6px; }
		a.btn:hover { background:#0069d9; }
		
		/* --- UPDATED STYLES FOR MARK APPLIED BUTTON --- */
		.btn-mark { 
			background:#495057; /* Darker, neutral background */
			color:#fff; 
			padding:8px 12px;
			border-radius:6px;
			cursor:pointer;
			border: none;
			font-weight: 500;
			transition: background 0.2s ease;
		}
		.btn-mark:hover { background:#343a40; } /* Subtle hover effect */
		.btn-mark:disabled {
			background: #e9ecef !important;
			color: #adb5bd !important;
			cursor: default;
		}
		/* --- END UPDATED STYLES --- */

		.header { display:flex; justify-content:space-between; align-items:center; margin-bottom: 12px; }
		.back { text-decoration:none; color:#007bff; font-weight: 500; }
		.actions { display:flex; gap:10px; align-items:center; }
		.download { background:#6c757d; color:#fff; padding:8px 14px; border-radius:6px; text-decoration:none; font-size:0.95em; }
		.download:hover { background:#5a6268; }
		.status-applied { opacity: 0.6; }
	</style>
	<script>
	  function markApplied(jobId, btn) {
		fetch('/mark-applied', {
		  method: 'POST',
		  headers: {'Content-Type': 'application/json'},
		  body: JSON.stringify({id: jobId})
		})
		.then(r => r.json())
		.then(data => {
		  if(data.success) {
			// Updated text to be more professional and use an icon (U+2705 is a white heavy check mark)
			btn.textContent = '✅ Applied'; 
			btn.disabled = true;
			btn.style.background = '#28a745'; // Use a success green color after application is marked
			btn.style.color = '#fff';
			btn.parentElement.parentElement.classList.add('status-applied');
		  } else {
			alert('Failed to update status');
		  }
		})
		.catch(() => alert('Error updating status'));
	  }
	</script>
 </head>
 <body>
   <div class="container">
	 <div class="header">
	   <h1>Jobs ({{.Total}} results)</h1>
	   <div class="actions">
		 <a class="download" href="/download-csv?{{.QueryString}}">⬇ Download CSV</a>
		 <a class="back" href="/filters">◀ Back to Filters</a>
	   </div>
	 </div>
	 <div style="margin-bottom:10px;">
	   {{if .Query}}<span class="pill">Search: {{.Query}}</span>{{end}}
	   {{if .Company}}<span class="pill">Company: {{.Company}}</span>{{end}}
	   {{if .Location}}<span class="pill">Location: {{.Location}}</span>{{end}}
	   {{if .Status}}<span class="pill">Status: {{.Status}}</span>{{end}}
	   {{range .Levels}}<span class="pill">{{.}}</span>{{end}}
	 </div>
	 <ul>
		{{range .Jobs}}
		<li {{if eq .Status "Applied"}}class="status-applied"{{end}}>
		   <div>
			  <div><strong>{{.Title}}</strong> — {{.Company}}</div>
			  <div class="meta">{{.Location}} • {{.Type}} • {{.DateAdded.Format "2006-01-02"}} • {{.Status}}</div>
		   </div>
		   <div>
			  <a class="btn" href="{{.URL}}" target="_blank">Open</a>
			  {{if eq .Status "Not Applied"}}
			  <button class="btn btn-mark" onclick="markApplied({{.ID}}, this)">Mark as Applied</button>
			  {{end}}
		   </div>
		</li>
		{{else}}
		<li>No jobs match your filters.</li>
		{{end}}
	 </ul>
   </div>
 </body>
 </html>
`
// -------------------- HELPERS --------------------

var levelRegex = regexp.MustCompile("(?i)(intern|new grad|new graduate|entry level|entry-level|junior|associate|apprentice|co-op|co op|coop|fellow)")

// deriveLevels returns canonical level labels found in a job title
func deriveLevels(title string) []string {
	title = strings.ToLower(title)
	uniqueLevels := make(map[string]struct{})
	levels := []string{}

	matches := levelRegex.FindAllString(title, -1)
	for _, match := range matches {
		switch match {
		case "intern":
			uniqueLevels["Intern"] = struct{}{}
		case "new grad", "new graduate":
			uniqueLevels["New Grad"] = struct{}{}
		case "entry level", "entry-level":
			uniqueLevels["Entry Level"] = struct{}{}
		case "junior":
			uniqueLevels["Junior"] = struct{}{}
		case "associate":
			uniqueLevels["Associate"] = struct{}{}
		case "apprentice":
			uniqueLevels["Apprentice"] = struct{}{}
		case "co-op", "co op", "coop":
			uniqueLevels["Co-op"] = struct{}{}
		case "fellow":
			uniqueLevels["Fellow"] = struct{}{}
		}
	}

	for level := range uniqueLevels {
		levels = append(levels, level)
	}
	sort.Strings(levels)
	return levels
}

// -------------------- HANDLERS (Authenticated) --------------------

// root handler redirects to login or filters
func rootHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["user"] == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Directs logged-in users to the main dashboard page
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// loginHandler handles GET (show form) and POST (process login)
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, loginHTML)
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		d, err := db.Connect()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		defer d.Close()
		valid, err := d.AuthenticateUser(username, password)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if !valid {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		session, _ := store.Get(r, "session")
		session.Values["user"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// registerHandler handles GET (show form) and POST (process registration)
func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, registerHTML)
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		d, err := db.Connect()
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		defer d.Close()
		err = d.CreateUser(username, password)
		if err != nil {
			http.Error(w, "Username already exists or invalid password", http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// logoutHandler clears the session
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Options.MaxAge = -1 // Expire the cookie
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusFound)
}

// filtersHandler shows the job filter page (authenticated)
func filtersHandler(w http.ResponseWriter, r *http.Request) {
	d, err := db.Connect()
	if err != nil {
		logger.Error("db connect: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Collect distinct titles and derive available levels
	rows, err := d.Conn.Query(`SELECT DISTINCT title FROM job_applications`)
	if err != nil {
		logger.Error("distinct titles: %v", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	levelSet := map[string]bool{}
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err == nil {
			for _, lv := range deriveLevels(title) {
				if lv != "" {
					levelSet[lv] = true
				}
			}
		}
	}
	rows.Close()
	levels := make([]string, 0, len(levelSet))
	for k := range levelSet {
		levels = append(levels, k)
	}
	sort.Strings(levels)

	// Collect distinct companies
	rows, err = d.Conn.Query(`SELECT DISTINCT company FROM job_applications ORDER BY company`)
	if err != nil {
		logger.Error("distinct companies: %v", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	var companies []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err == nil && c != "" {
			companies = append(companies, c)
		}
	}
	rows.Close()

	// Collect distinct locations
	rows, err = d.Conn.Query(`SELECT DISTINCT location FROM job_applications ORDER BY location`)
	if err != nil {
		logger.Error("distinct locations: %v", err)
		http.Error(w, "Query error", http.StatusInternalServerError)
		return
	}
	var locations []string
	for rows.Next() {
		var l string
		if err := rows.Scan(&l); err == nil && l != "" {
			locations = append(locations, l)
		}
	}
	rows.Close()

	// Get username for the template
	session, _ := store.Get(r, "session")
	user := fmt.Sprintf("%v", session.Values["user"])

	// Render the template
	lt := template.Must(template.New("landing").Parse(landingHTML))
	data := struct {
		Levels    []string
		Companies []string
		Locations []string
		User      string
	}{Levels: levels, Companies: companies, Locations: locations, User: user}
	if err := lt.Execute(w, data); err != nil {
		logger.Error("landing template: %v", err)
	}
}

// dashboardHandler shows the job stats (authenticated)
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	d, err := db.Connect()
	if err != nil {
		logger.Error("db connect: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// 1. Calculate stats (TOTAL, NOT APPLIED, APPLIED) across ALL jobs
	var totalCount, notAppliedCount, appliedCount int
	err = d.Conn.QueryRow(`
		SELECT COUNT(*), 
			SUM(CASE WHEN status = 'Not Applied' THEN 1 ELSE 0 END), 
			SUM(CASE WHEN status = 'Applied' THEN 1 ELSE 0 END) 
		FROM job_applications`).Scan(&totalCount, &notAppliedCount, &appliedCount)
	if err != nil && err != sql.ErrNoRows {
		logger.Error("query job stats: %v", err)
		http.Error(w, "Query error for stats", http.StatusInternalServerError)
		return
	}
	
	// 2. Query for JOBS, but filter to show ONLY 'Applied' jobs on the dashboard page
	rows, err := d.Conn.Query(`
		SELECT id, title, company, location, type, url, date_added, status 
		FROM job_applications
		WHERE status = 'Applied' 
		ORDER BY date_added DESC`)
	if err != nil {
		logger.Error("query applied jobs: %v", err)
		http.Error(w, "Query error for jobs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		var typ string
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &typ, &job.URL, &job.DateAdded, &job.Status); err != nil {
			logger.Error("scan row: %v", err)
			continue
		}
		job.Type = typ
		jobs = append(jobs, job)
	}

	session, _ := store.Get(r, "session")
	user := fmt.Sprintf("%v", session.Values["user"])

	data := PageData{
		Jobs:       jobs,
		TotalJobs:  totalCount,
		NotApplied: notAppliedCount,
		Applied:    appliedCount,
		User:       user,
	}

	// template with a small helper to create CSS class names
	tmpl := template.Must(template.New("dashboard").Funcs(template.FuncMap{"lower": func(s string) string { return strings.ToLower(strings.ReplaceAll(s, " ", "-")) }}).Parse(dashboardHTML))

	if err := tmpl.Execute(w, data); err != nil {
		logger.Error("template execute: %v", err)
	}
}

// resultsHandler shows filtered job results (authenticated, includes pagination)
func resultsHandler(w http.ResponseWriter, r *http.Request) {
	// Query parameters
	selLevels := r.URL.Query()["level"]
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	status := r.URL.Query().Get("status")
	company := r.URL.Query().Get("company")
	location := r.URL.Query().Get("location")

	d, err := db.Connect()
	if err != nil {
		logger.Error("db connect: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Build WHERE clause
	var clauses []string
	var args []interface{}

	// Status filter
	if status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, status)
	}
	// Company filter
	if company != "" {
		clauses = append(clauses, fmt.Sprintf("company = $%d", len(args)+1))
		args = append(args, company)
	}
	// Location filter
	if location != "" {
		clauses = append(clauses, fmt.Sprintf("location = $%d", len(args)+1))
		args = append(args, location)
	}
	// Query string over title
	if q != "" {
		// Use LIKE on SQLite, with COLLATE NOCASE for case-insensitivity
		clauses = append(clauses, fmt.Sprintf("title LIKE $%d COLLATE NOCASE", len(args)+1))
		args = append(args, "%"+q+"%")
	}
	// Levels mapped to title keywords
	var levelPatterns []string
	for _, lv := range selLevels {
		switch lv {
		case "Intern":
			levelPatterns = append(levelPatterns, "intern")
		case "New Grad":
			levelPatterns = append(levelPatterns, "new grad", "new graduate")
		case "Entry Level":
			levelPatterns = append(levelPatterns, "entry level", "entry-level")
		case "Junior":
			levelPatterns = append(levelPatterns, "junior")
		case "Associate":
			levelPatterns = append(levelPatterns, "associate")
		case "Apprentice":
			levelPatterns = append(levelPatterns, "apprentice")
		case "Fellow":
			levelPatterns = append(levelPatterns, "fellow")
		case "Co-op":
			levelPatterns = append(levelPatterns, "co-op", "co op", "coop")
		}
	}

	if len(levelPatterns) > 0 {
		// Build (title LIKE $x OR title LIKE $y ...) group
		var parts []string
		for _, pat := range levelPatterns {
			parts = append(parts, fmt.Sprintf("title LIKE $%d COLLATE NOCASE", len(args)+1))
			args = append(args, "%"+pat+"%")
		}
		clauses = append(clauses, "("+strings.Join(parts, " OR ")+")")
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	var jobs []Job

	// Pagination: page & page_size
	page := 1
	pageSize := 20
	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil {
			if v < 1 {
				v = 1
			}
			if v > 200 {
				v = 200
			}
			pageSize = v
		}
	}

	// Get total count for this query
	var total int
	countQ := "SELECT COUNT(*) FROM job_applications" + where
	if err := d.Conn.QueryRow(countQ, args...).Scan(&total); err != nil {
		logger.Error("count query: %v", err)
	}

	// Apply limit/offset
	offset := (page - 1) * pageSize
	limitIdx := len(args) + 1
	offsetIdx := len(args) + 2

	dataQ := fmt.Sprintf(
		"SELECT id, title, company, location, type, url, date_added, status FROM job_applications%s ORDER BY date_added DESC LIMIT $%d OFFSET $%d",
		where, limitIdx, offsetIdx,
	)
	argsData := append([]interface{}{}, args...)
	argsData = append(argsData, pageSize, offset)

	rows, err := d.Conn.Query(dataQ, argsData...)
	if err != nil {
		logger.Error("list query: %v", err)
		http.Error(w, `{"error":"list query"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var job Job
		var typ string
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &typ, &job.URL, &job.DateAdded, &job.Status); err != nil {
			logger.Error("scan row: %v", err)
			continue
		}
		job.Type = typ
		jobs = append(jobs, job)
	}

	// Struct used by the template (simplified to match the older UI format)
	rt := template.Must(template.New("results").Funcs(template.FuncMap{"eq": func(a, b interface{}) bool { return a == b }}).Parse(resultsHTML))
	data := struct {
		Jobs        []Job
		Levels      []string
		Query       string
		Company     string
		Location    string
		Status      string
		Total       int
		QueryString string
	}{
		Jobs: jobs, Levels: selLevels, Query: q, Company: company, Location: location,
		Status: status, Total: total, QueryString: r.URL.RawQuery,
	}
	if err := rt.Execute(w, data); err != nil {
		logger.Error("results template: %v", err)
	}
}

// downloadCSVHandler exports filtered job results as CSV (authenticated)
func downloadCSVHandler(w http.ResponseWriter, r *http.Request) {
	selLevels := r.URL.Query()["level"]
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	status := r.URL.Query().Get("status")
	company := r.URL.Query().Get("company")
	location := r.URL.Query().Get("location")

	d, err := db.Connect()
	if err != nil {
		logger.Error("db connect: %v", err)
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Build WHERE clause (same logic as /results but without pagination/user)
	var clauses []string
	var args []interface{}

	if status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, status)
	}
	if company != "" {
		clauses = append(clauses, fmt.Sprintf("company = $%d", len(args)+1))
		args = append(args, company)
	}
	if location != "" {
		clauses = append(clauses, fmt.Sprintf("location = $%d", len(args)+1))
		args = append(args, location)
	}
	if q != "" {
		clauses = append(clauses, fmt.Sprintf("title LIKE $%d COLLATE NOCASE", len(args)+1))
		args = append(args, "%"+q+"%")
	}
	var levelPatterns []string
	for _, lv := range selLevels {
		switch lv {
		case "Intern":
			levelPatterns = append(levelPatterns, "intern")
		case "New Grad":
			levelPatterns = append(levelPatterns, "new grad", "new graduate")
		case "Entry Level":
			levelPatterns = append(levelPatterns, "entry level", "entry-level")
		case "Junior":
			levelPatterns = append(levelPatterns, "junior")
		case "Associate":
			levelPatterns = append(levelPatterns, "associate")
		case "Apprentice":
			levelPatterns = append(levelPatterns, "apprentice")
		case "Fellow":
			levelPatterns = append(levelPatterns, "fellow")
		case "Co-op":
			levelPatterns = append(levelPatterns, "co-op", "co op", "coop")
		}
	}
	if len(levelPatterns) > 0 {
		var parts []string
		for _, pat := range levelPatterns {
			parts = append(parts, fmt.Sprintf("title LIKE $%d COLLATE NOCASE", len(args)+1))
			args = append(args, "%"+pat+"%")
		}
		clauses = append(clauses, "("+strings.Join(parts, " OR ")+")")
	}

	where := ""
	if len(clauses) > 0 {
		where = " WHERE " + strings.Join(clauses, " AND ")
	}

	// Select all fields for CSV
	dataQ := fmt.Sprintf(
		"SELECT title, company, location, type, url, date_added, status FROM job_applications%s ORDER BY date_added DESC",
		where,
	)
	rows, err := d.Conn.Query(dataQ, args...)
	if err != nil {
		logger.Error("list query: %v", err)
		http.Error(w, `{"error":"list query"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=jobs.csv")
	w.Write([]byte("Title,Company,Location,Type,URL,Date Added,Status\n"))

	// CSV escaping helper
	escape := func(s string) string {
		if strings.Contains(s, ",") || strings.Contains(s, "\"") {
			return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
		}
		return s
	}

	for rows.Next() {
		var title, company, loc, typ, url, dateAdded, st string
		if err := rows.Scan(&title, &company, &loc, &typ, &url, &dateAdded, &st); err != nil {
			continue
		}
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n", escape(title), escape(company), escape(loc), escape(typ), escape(url), escape(dateAdded), escape(st))
		w.Write([]byte(line))
	}
}

// markAppliedHandler updates job status via POST (authenticated)
func markAppliedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false}`, http.StatusBadRequest)
		return
	}
	d, err := db.Connect()
	if err != nil {
		logger.Error("db connect: %v", err)
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	defer d.Close()

	_, err = d.Conn.Exec(`UPDATE job_applications SET status = 'Applied' WHERE id = $1`, req.ID)
	if err != nil {
		logger.Error("update status: %v", err)
		http.Error(w, `{"success":false}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true}`))
}

// -------------------- MAIN --------------------

func main() {
	// Flag parsing remains (if needed)
	port := flag.String("port", "8080", "port to serve on")
	flag.Parse()

	// Connect to DB and ensure schemas are created
	d, err := db.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to DB: %v", err)
	}
	defer d.Close()

	// Create job applications table
	if err := d.CreateSchema(); err != nil {
		logger.Fatal("Failed to create job_applications schema: %v", err)
	}

	// Create users table
	if err := d.CreateUserSchema(); err != nil {
		logger.Fatal("Failed to create users schema: %v", err)
	}

	// All routes are now based on the authenticated logic
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Protected routes (UI from main (1).go, logic from main.go)
	http.HandleFunc("/filters", AuthRequired(filtersHandler))
	http.HandleFunc("/dashboard", AuthRequired(dashboardHandler))
	http.HandleFunc("/results", AuthRequired(resultsHandler))
	http.HandleFunc("/download-csv", AuthRequired(downloadCSVHandler))
	http.HandleFunc("/mark-applied", AuthRequired(markAppliedHandler))
	
	logger.Info("Listening on http://localhost:%s", *port)
	if err := http.ListenAndServe(":"+*port, nil); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}