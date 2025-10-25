package main

import (
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

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/logger"
)

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

type PageData struct {
	Jobs       []Job
	TotalJobs  int
	NotApplied int
	Applied    int
}

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
    <h1>🚀 Job Application Dashboard</h1>
    
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
		.btn-mark { background:#28a745; cursor:pointer; }
		.btn-mark:hover { background:#218838; }
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
			btn.textContent = '✓ Applied';
			btn.disabled = true;
			btn.style.background = '#6c757d';
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
		 <a class="back" href="/">◀ Back</a>
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
			  <button class="btn btn-mark" onclick="markApplied({{.ID}}, this)">Mark Applied</button>
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

// deriveLevels returns canonical level labels found in a job title
func deriveLevels(title string) []string {
	t := strings.ToLower(title)
	var out []string
	add := func(s string) { out = append(out, s) }
	// Canonical buckets
	if matched, _ := regexp.MatchString(`\bintern(ship)?\b`, t); matched {
		add("Intern")
	}
	if strings.Contains(t, "new grad") || strings.Contains(t, "new graduate") {
		add("New Grad")
	}
	if strings.Contains(t, "entry level") || strings.Contains(t, "entry-level") {
		add("Entry Level")
	}
	if strings.Contains(t, "junior") {
		add("Junior")
	}
	if strings.Contains(t, "associate") {
		add("Associate")
	}
	if strings.Contains(t, "apprentice") {
		add("Apprentice")
	}
	if strings.Contains(t, "fellow") {
		add("Fellow")
	}
	if strings.Contains(t, "co-op") || strings.Contains(t, "co op") || strings.Contains(t, "coop") {
		add("Co-op")
	}
	// If none matched but looks generic early career, classify as Entry Level
	if len(out) == 0 {
		// Heuristic: contains engineer/developer/analyst without senior keywords
		if matched, _ := regexp.MatchString(`\b(engineer|developer|analyst|specialist|coordinator)\b`, t); matched {
			if ok, _ := regexp.MatchString(`\b(senior|staff|principal|lead|manager|director|architect|head|chief|vp)\b`, t); !ok {
				add("Entry Level")
			}
		}
	}
	// dedupe
	if len(out) > 1 {
		seen := map[string]bool{}
		uniq := []string{}
		for _, v := range out {
			if !seen[v] {
				seen[v] = true
				uniq = append(uniq, v)
			}
		}
		out = uniq
	}
	return out
}

func main() {
	port := flag.String("port", "8080", "Port to run dashboard on")
	flag.Parse()

	// Init logger from env
	logger.InitFromEnv()

	logger.Info("Starting Job Dashboard Server on port %s", *port)

	tmpl := template.Must(template.New("dashboard").Funcs(template.FuncMap{
		"lower": func(s string) string {
			if s == "Not Applied" {
				return "not-applied"
			}
			return "applied"
		},
	}).Parse(dashboardHTML))

	// Landing page with dynamic levels, companies, and locations
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			if err := rows.Scan(&title); err != nil {
				continue
			}
			for _, lv := range deriveLevels(title) {
				if lv != "" {
					levelSet[lv] = true
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
			var loc string
			if err := rows.Scan(&loc); err == nil && loc != "" {
				locations = append(locations, loc)
			}
		}
		rows.Close()

		lt := template.Must(template.New("landing").Parse(landingHTML))
		data := struct {
			Levels    []string
			Companies []string
			Locations []string
		}{Levels: levels, Companies: companies, Locations: locations}
		if err := lt.Execute(w, data); err != nil {
			logger.Error("landing template: %v", err)
		}
	})

	// Results page
	http.HandleFunc("/results", func(w http.ResponseWriter, r *http.Request) {
		// Parse inputs
		selLevels := r.URL.Query()["level"] // can be multiple
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

		// Build dynamic WHERE clause
		var clauses []string
		var args []interface{}
		// Status
		if status != "" {
			clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)+1))
			args = append(args, status)
		}
		// Company
		if company != "" {
			clauses = append(clauses, fmt.Sprintf("company = $%d", len(args)+1))
			args = append(args, company)
		}
		// Location
		if location != "" {
			clauses = append(clauses, fmt.Sprintf("location = $%d", len(args)+1))
			args = append(args, location)
		}
		// Query string over title
		if q != "" {
			clauses = append(clauses, fmt.Sprintf("title ILIKE $%d", len(args)+1))
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
			// Build (title ILIKE $x OR title ILIKE $y ...) group
			var parts []string
			for _, pat := range levelPatterns {
				parts = append(parts, fmt.Sprintf("title ILIKE $%d", len(args)+1))
				args = append(args, "%"+pat+"%")
			}
			clauses = append(clauses, "("+strings.Join(parts, " OR ")+")")
		}

		where := ""
		if len(clauses) > 0 {
			where = " WHERE " + strings.Join(clauses, " AND ")
		}

		query := "SELECT id, title, company, location, type, url, date_added, status FROM job_applications" + where + " ORDER BY date_added DESC LIMIT 500"
		rows, err := d.Conn.Query(query, args...)
		if err != nil {
			logger.Error("results query: %v", err)
			http.Error(w, "Query error", http.StatusInternalServerError)
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

		rt := template.Must(template.New("results").Parse(resultsHTML))
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
			Jobs:        jobs,
			Levels:      selLevels,
			Query:       q,
			Company:     company,
			Location:    location,
			Status:      status,
			Total:       len(jobs),
			QueryString: r.URL.RawQuery,
		}

		if err := rt.Execute(w, data); err != nil {
			logger.Error("results template: %v", err)
		}
	})

	// CSV download for filtered results
	http.HandleFunc("/download-csv", func(w http.ResponseWriter, r *http.Request) {
		// Parse same filters as /results
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

		// Build WHERE clause (same logic as /results)
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
			clauses = append(clauses, fmt.Sprintf("title ILIKE $%d", len(args)+1))
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
				parts = append(parts, fmt.Sprintf("title ILIKE $%d", len(args)+1))
				args = append(args, "%"+pat+"%")
			}
			clauses = append(clauses, "("+strings.Join(parts, " OR ")+")")
		}

		where := ""
		if len(clauses) > 0 {
			where = " WHERE " + strings.Join(clauses, " AND ")
		}

		query := "SELECT title, company, location, type, url, date_added, status FROM job_applications" + where + " ORDER BY date_added DESC LIMIT 500"
		rows, err := d.Conn.Query(query, args...)
		if err != nil {
			logger.Error("csv query: %v", err)
			http.Error(w, "Query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=jobs.csv")
		w.Write([]byte("Title,Company,Location,Type,URL,Date Added,Status\n"))

		for rows.Next() {
			var title, company, loc, typ, url, dateAdded, st string
			if err := rows.Scan(&title, &company, &loc, &typ, &url, &dateAdded, &st); err != nil {
				continue
			}
			// Simple CSV escaping: quote if contains comma or quote
			escape := func(s string) string {
				if strings.Contains(s, ",") || strings.Contains(s, "\"") {
					return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
				}
				return s
			}
			line := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s\n",
				escape(title), escape(company), escape(loc), escape(typ), escape(url), escape(dateAdded), escape(st))
			w.Write([]byte(line))
		}
	})

	// POST endpoint to mark a job as Applied
	http.HandleFunc("/mark-applied", func(w http.ResponseWriter, r *http.Request) {
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
	})

	// Preserve the original dashboard at /dashboard
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		d, err := db.Connect()
		if err != nil {
			logger.Error("db connect: %v", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		defer d.Close()

		rows, err := d.Conn.Query(`
			SELECT id, title, company, location, type, url, date_added, status 
			FROM job_applications 
			ORDER BY date_added DESC
		`)
		if err != nil {
			logger.Error("query jobs: %v", err)
			http.Error(w, "Query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var jobs []Job
		notApplied := 0
		applied := 0

		for rows.Next() {
			var job Job
			var typ string
			err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &typ, &job.URL, &job.DateAdded, &job.Status)
			if err != nil {
				logger.Error("scan row: %v", err)
				continue
			}
			job.Type = typ
			jobs = append(jobs, job)

			if job.Status == "Not Applied" {
				notApplied++
			} else {
				applied++
			}
		}

		data := PageData{
			Jobs:       jobs,
			TotalJobs:  len(jobs),
			NotApplied: notApplied,
			Applied:    applied,
		}

		if err := tmpl.Execute(w, data); err != nil {
			logger.Error("template execute: %v", err)
		}
	})

	// JSON API: /api/jobs?page=1&page_size=50&status=Applied|Not%20Applied
	http.HandleFunc("/api/jobs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Parse query params
		page := 1
		pageSize := 50
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
		statusFilter := r.URL.Query().Get("status")

		d, err := db.Connect()
		if err != nil {
			logger.Error("db connect: %v", err)
			http.Error(w, `{"error":"db connect"}`, http.StatusInternalServerError)
			return
		}
		defer d.Close()

		// Count total
		where := ""
		args := []interface{}{}
		if statusFilter != "" {
			where = " WHERE status = $1"
			args = append(args, statusFilter)
		}

		var total int
		countQ := "SELECT COUNT(*) FROM job_applications" + where
		if err := d.Conn.QueryRow(countQ, args...).Scan(&total); err != nil {
			logger.Error("count query: %v", err)
			http.Error(w, `{"error":"count query"}`, http.StatusInternalServerError)
			return
		}

		// Fetch page
		offset := (page - 1) * pageSize
		// Build placeholders for limit and offset based on args length
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
		if err := rows.Err(); err != nil {
			logger.Error("rows err: %v", err)
			http.Error(w, `{"error":"rows"}`, http.StatusInternalServerError)
			return
		}

		totalPages := (total + pageSize - 1) / pageSize
		resp := struct {
			Page       int   `json:"page"`
			PageSize   int   `json:"page_size"`
			Total      int   `json:"total"`
			TotalPages int   `json:"total_pages"`
			Items      []Job `json:"items"`
		}{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
			Items:      jobs,
		}

		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(resp); err != nil {
			logger.Error("encode json: %v", err)
		}
	})

	logger.Info("Dashboard available at http://localhost:%s", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), nil); err != nil {
		logger.Fatal("server error: %v", err)
	}
}
