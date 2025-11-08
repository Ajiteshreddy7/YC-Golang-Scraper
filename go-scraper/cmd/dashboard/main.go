package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/logger"
	"github.com/gorilla/sessions"
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

var loginHTML = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<title>YC Job Scraper — Intelligent Job Tracker</title>
	<style>
		:root{--bg:#f4f6f8;--card:#fff;--accent:#0b66c3;--muted:#6b7280;--border:#e2e8f0}
		*{box-sizing:border-box}
		html,body{height:100%;margin:0;font-family:Inter,ui-sans-serif,system-ui,-apple-system,"Segoe UI",Roboto,"Helvetica Neue",Arial}
		body{background:#f8fafc;display:flex;align-items:center;justify-content:center;padding:20px;min-height:100vh}
		.container{max-width:440px;width:100%}
		.card{background:var(--card);border-radius:16px;box-shadow:0 20px 40px rgba(0,0,0,0.15);padding:40px;border:1px solid rgba(255,255,255,0.2)}
		.brand{text-align:center;margin-bottom:32px}
		.logo{width:48px;height:48px;border-radius:12px;background:linear-gradient(135deg,var(--accent),#2aa6ff);display:inline-block;margin-bottom:16px;position:relative}
		.logo::after{content:"🚀";position:absolute;top:50%;left:50%;transform:translate(-50%,-50%);font-size:20px}
		h1{margin:0 0 8px;font-size:24px;color:#111827;font-weight:700;text-align:center}
		p.lead{margin:0 0 32px;color:var(--muted);font-size:14px;text-align:center;line-height:1.5}
		.form-group{margin-bottom:20px}
		label{display:block;font-size:14px;color:#374151;margin-bottom:8px;font-weight:500}
		input[type=text],input[type=password]{width:100%;padding:14px 16px;border:2px solid var(--border);border-radius:10px;font-size:15px;transition:border-color 0.2s ease,box-shadow 0.2s ease;background:#fff}
		input[type=text]:focus,input[type=password]:focus{outline:none;border-color:var(--accent);box-shadow:0 0 0 3px rgba(11,102,195,0.1)}
		input[type=text]:hover,input[type=password]:hover{border-color:#cbd5e0}
		.remember-me{margin:16px 0 24px;display:flex;align-items:center;gap:8px}
		.remember-me input[type=checkbox]{width:auto;margin:0}
		.remember-me label{margin:0;font-size:14px;color:var(--muted);font-weight:400}
		.btn-group{display:flex;flex-direction:column;gap:12px;margin-top:8px}
		button.primary{background:var(--accent);color:#fff;border:none;padding:14px 20px;border-radius:10px;font-size:15px;font-weight:600;cursor:pointer;width:100%;transition:background-color 0.2s ease,transform 0.1s ease}
		button.primary:hover{background:#0a5aa8;transform:translateY(-1px)}
		button.primary:active{transform:translateY(0)}
		a.link{color:var(--accent);text-decoration:none;font-weight:500;text-align:center;padding:12px;border-radius:8px;transition:background-color 0.2s ease}
		a.link:hover{background:#f8fafc}
		.footer{margin-top:24px;padding-top:20px;border-top:1px solid #f1f5f9;font-size:13px;color:var(--muted);text-align:center;line-height:1.6}
		.footer a{color:var(--accent);text-decoration:none}
		.footer a:hover{text-decoration:underline}
		.demo-info{margin-top:12px;padding:12px;background:#f8fafc;border:1px solid #e2e8f0;border-radius:8px;font-size:13px;color:#4a5568}
		.demo-info strong{color:#2d3748}
		/* toast */
		.toast{position:fixed;right:20px;top:20px;background:#fff;border-left:4px solid #f87171;padding:12px 16px;border-radius:8px;box-shadow:0 10px 25px rgba(0,0,0,0.15);display:none;z-index:1000}
		.toast.show{display:block}
		@media (max-width: 480px) {
			body{padding:16px}
			.card{padding:32px 24px}
			h1{font-size:22px}
		}
	</style>
</head>
<body>
	<div class="container">
		<div class="card">
			<div class="brand">
				<div class="logo"></div>
				<h1>YC Job Scraper</h1>
				<p class="lead">Intelligent early-career job tracker — Sign in to manage your applications</p>
			</div>

			<form method="POST" action="/login" id="loginForm">
				<div class="form-group">
					<label for="username">Username</label>
					<input id="username" name="username" type="text" autocomplete="username" placeholder="Enter your username" required />
				</div>

				<div class="form-group">
					<label for="password">Password</label>
					<input id="password" name="password" type="password" autocomplete="current-password" placeholder="Enter your password" required />
				</div>

				<div class="remember-me">
					<input type="checkbox" id="remember" name="remember" />
					<label for="remember">Remember me for 30 days</label>
				</div>

				<div class="btn-group">
					<button class="primary" type="submit">Sign in</button>
					<a class="link" href="#" onclick="alert('Registration coming soon!')">Create new account</a>
				</div>
			</form>

			<div class="footer">
				By signing in you agree to our <a href="#">Terms</a> and <a href="#">Privacy Policy</a>.
				<div class="demo-info">
					<strong>Demo Credentials:</strong> admin / password123
				</div>
			</div>
		</div>
	</div>

	<div id="toast" class="toast" role="alert" aria-live="assertive"></div>

	<script>
		(function(){
			// Show toast if server rendered an error
			var toast = document.getElementById('toast');
			var err = '{{.Error}}';
			if(err && err !== ''){
				toast.textContent = err;
				toast.className = 'toast show';
				setTimeout(function(){ toast.className = 'toast'; }, 4500);
			}
		})();
	</script>
</body>
</html>`

var registerHTML = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<title>Create account — Go Scrape</title>
	<style>
		:root{--bg:#f4f6f8;--card:#fff;--accent:#0b66c3;--muted:#6b7280}
		html,body{height:100%;margin:0;font-family:Inter,ui-sans-serif,system-ui,-apple-system,"Segoe UI",Roboto,"Helvetica Neue",Arial}
		body{background:linear-gradient(180deg,#f8fafc 0%,var(--bg) 100%);display:flex;align-items:center;justify-content:center;padding:32px}
		.container{max-width:460px;width:100%}
		.card{background:var(--card);border-radius:12px;box-shadow:0 10px 30px rgba(2,6,23,0.08);padding:28px}
		.brand{display:flex;align-items:center;gap:12px;margin-bottom:18px}
		.logo{width:44px;height:44px;border-radius:10px;background:linear-gradient(135deg,var(--accent),#2aa6ff);display:inline-block}
		h1{margin:0;font-size:20px;color:#0f172a}
		p.lead{margin:6px 0 18px;color:var(--muted);font-size:13px}
		label{display:block;font-size:13px;color:#111827;margin-bottom:6px}
		input[type=text],input[type=password]{width:100%;padding:12px 14px;border:1px solid #e6e9ee;border-radius:8px;margin-bottom:12px;font-size:14px}
		.actions{display:flex;align-items:center;justify-content:space-between;margin-top:6px}
		button.primary{background:var(--accent);color:#fff;border:none;padding:10px 14px;border-radius:8px;font-weight:600;cursor:pointer}
		a.link{color:var(--accent);text-decoration:none;font-weight:600}
		.footer{margin-top:14px;font-size:13px;color:var(--muted);text-align:center}
		.toast{position:fixed;right:20px;top:20px;background:#fff;border-left:4px solid #f87171;padding:12px 16px;border-radius:6px;box-shadow:0 6px 24px rgba(2,6,23,0.08);display:none}
		.toast.show{display:block}
	</style>
</head>
<body>
	<div class="container">
		<div class="card">
			<div class="brand">
				<span class="logo" aria-hidden="true"></span>
				<div>
					<h1>Create your Go Scrape account</h1>
					<div class="lead">Jobs at One sight — sign up to track applications</div>
				</div>
			</div>

			<form method="POST" action="/register" id="registerForm">
				<label for="username">Username</label>
				<input id="username" name="username" type="text" autocomplete="username" required />

				<label for="password">Password</label>
				<input id="password" name="password" type="password" autocomplete="new-password" required />

				<label for="confirm">Confirm Password</label>
				<input id="confirm" name="confirm" type="password" autocomplete="new-password" required />

				<div style="margin-top:18px;display:flex;gap:8px">
					<button class="primary" type="submit">Create account</button>
					<a class="link" href="/login" style="align-self:center">Sign in</a>
				</div>
			</form>

			<div class="footer">By creating an account you agree to our <a href="#">Terms</a> and <a href="#">Privacy Policy</a>.</div>
		</div>
	</div>

	<div id="toast" class="toast" role="alert" aria-live="assertive"></div>

	<script>
		(function(){
			var toast = document.getElementById('toast');
			var err = '{{.Error}}';
			if(err && err !== ''){
				toast.textContent = err;
				toast.className = 'toast show';
				setTimeout(function(){ toast.className = 'toast'; }, 4500);
			}
		})();
	</script>
</body>
</html>`

// -------------------- UI TEMPLATES (from main (1).go) --------------------

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
        <div style="font-size: 0.95em; color: #6c757d;">Logged in as <strong>{{.User}}</strong> — <a href="/results?status=Not%20Applied">Dashboard</a> | <a href="/logout">Logout</a></div>
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
				<option value="Not Applied" selected>Not Applied</option>
				<option value="Applied">Applied</option>
				<option value="">All Statuses</option>
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
	<title>Job Dashboard</title>
	<style>
		body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; margin: 40px; background: #f8f9fa; }
		.container { max-width: 1100px; margin: 0 auto; }
		h1 { color: #343a40; }
		/* Dashboard Stats Styling */
		.stats {
			display: flex;
			justify-content: space-around;
			margin: 20px 0 30px 0;
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
		/* Filter Pills and Job List Styling */
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
	   <h1>Job Dashboard ({{.Total}} jobs found)</h1>
	   <div class="actions">
		 <a class="download" href="/download-csv?{{.QueryString}}">⬇ Download CSV</a>
		 <a class="back" href="/filters">◀ Back to Filters</a>
	   </div>
	 </div>

	 <!-- Dashboard Stats Section -->
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
	// For debugging: Add a simple health check response if requested
	if r.URL.Query().Get("health") == "check" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("OK - YC Job Tracker is running"))
		return
	}

	session, _ := store.Get(r, "session")
	if session.Values["user"] == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	// Directs logged-in users to the filters page as the main landing page
	http.Redirect(w, r, "/filters", http.StatusFound)
}

// loginHandler handles GET (show form) and POST (process login)
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Use a small template to render the login page with an optional error message
	tmpl := template.Must(template.New("login").Parse(loginHTML))
	switch r.Method {
	case http.MethodGet:
		// Render with no error
		_ = tmpl.Execute(w, map[string]string{"Error": ""})
		return
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		d, err := db.Connect()
		if err != nil {
			logger.Error("db connect: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = tmpl.Execute(w, map[string]string{"Error": "Internal server error"})
			return
		}
		defer d.Close()
		valid, err := d.AuthenticateUser(username, password)
		if err != nil {
			logger.Error("auth check: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = tmpl.Execute(w, map[string]string{"Error": "Internal server error"})
			return
		}
		if !valid {
			// Render the login page with a toast-style error (no redirect)
			w.WriteHeader(http.StatusUnauthorized)
			_ = tmpl.Execute(w, map[string]string{"Error": "Invalid username or password"})
			return
		}
		session, _ := store.Get(r, "session")
		session.Values["user"] = username
		session.Save(r, w)
		http.Redirect(w, r, "/filters", http.StatusFound)
		return
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// registerHandler handles GET (show form) and POST (process registration)
func registerHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("register").Parse(registerHTML))
	switch r.Method {
	case http.MethodGet:
		_ = tmpl.Execute(w, map[string]string{"Error": ""})
	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		confirm := strings.TrimSpace(r.FormValue("confirm"))

		if username == "" || password == "" {
			w.WriteHeader(http.StatusBadRequest)
			_ = tmpl.Execute(w, map[string]string{"Error": "Username and password are required"})
			return
		}
		if password != confirm {
			w.WriteHeader(http.StatusBadRequest)
			_ = tmpl.Execute(w, map[string]string{"Error": "Passwords do not match"})
			return
		}

		d, err := db.Connect()
		if err != nil {
			logger.Error("db connect: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = tmpl.Execute(w, map[string]string{"Error": "Internal server error"})
			return
		}
		defer d.Close()

		err = d.CreateUser(username, password)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = tmpl.Execute(w, map[string]string{"Error": "Username already exists or invalid password"})
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

// dashboardHandler redirects to the results page (which now serves as the main dashboard)
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Redirect to results page with default filters (Not Applied status)
	http.Redirect(w, r, "/results?status=Not%20Applied", http.StatusFound)
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

	// First, get overall job statistics for the dashboard counters
	var totalCount, notAppliedCount, appliedCount int
	err = d.Conn.QueryRow(`
		SELECT COUNT(*), 
			COALESCE(SUM(CASE WHEN status = 'Not Applied' THEN 1 ELSE 0 END), 0), 
			COALESCE(SUM(CASE WHEN status = 'Applied' THEN 1 ELSE 0 END), 0) 
		FROM job_applications`).Scan(&totalCount, &notAppliedCount, &appliedCount)
	if err != nil {
		logger.Error("query job stats: %v", err)
		http.Error(w, "Query error for stats", http.StatusInternalServerError)
		return
	}

	// Build WHERE clause
	var clauses []string
	var args []interface{}

	// Status filter - default to "Not Applied" unless explicitly set
	if status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, status)
	} else {
		// By default, only show "Not Applied" jobs
		clauses = append(clauses, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, "Not Applied")
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

	// Struct used by the template (enhanced with job statistics)
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
		TotalJobs   int
		NotApplied  int
		Applied     int
	}{
		Jobs: jobs, Levels: selLevels, Query: q, Company: company, Location: location,
		Status: status, Total: total, QueryString: r.URL.RawQuery,
		TotalJobs: totalCount, NotApplied: notAppliedCount, Applied: appliedCount,
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

// initAdminHandler creates admin user if it doesn't exist (for Render deployment)
func initAdminHandler(w http.ResponseWriter, r *http.Request) {
	d, err := db.Connect()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Check if admin user already exists
	_, _, _, err = d.GetUserByUsername("admin")
	if err == nil {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
		<!DOCTYPE html>
		<html><head><title>Admin Already Exists</title></head>
		<body style="font-family: Arial; text-align: center; padding: 50px;">
			<h2>✅ Admin User Already Exists</h2>
			<p>The admin user is already configured.</p>
			<p><strong>Username:</strong> admin</p>
			<p><strong>Password:</strong> password123</p>
			<a href="/login" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Go to Login</a>
		</body></html>`))
		return
	}

	// Create admin user
	err = d.CreateUser("admin", "password123")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create admin user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`
	<!DOCTYPE html>
	<html><head><title>Admin User Created</title></head>
	<body style="font-family: Arial; text-align: center; padding: 50px;">
		<h2>✅ Admin User Created Successfully!</h2>
		<p>Your admin account has been set up.</p>
		<p><strong>Username:</strong> admin</p>
		<p><strong>Password:</strong> password123</p>
		<p><em>Please change the password after first login</em></p>
		<a href="/login" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Go to Login</a>
	</body></html>`))
}

// JobImport represents the structure for importing jobs
type JobImport struct {
	Title     string    `json:"title"`
	Company   string    `json:"company"`
	Location  string    `json:"location"`
	Type      string    `json:"type"`
	URL       string    `json:"url"`
	DateAdded time.Time `json:"date_added"`
	Status    string    `json:"status"`
}

// importJobsHandler imports jobs from a JSON URL (for Render deployment)
func importJobsHandler(w http.ResponseWriter, r *http.Request) {
	jsonURL := r.URL.Query().Get("url")
	if jsonURL == "" {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
		<!DOCTYPE html>
		<html><head><title>Import Jobs</title></head>
		<body style="font-family: Arial; text-align: center; padding: 50px;">
			<h2>📥 Import Jobs to Render</h2>
			<p>Enter the URL of your jobs JSON file:</p>
			<form method="GET">
				<input type="url" name="url" placeholder="https://example.com/jobs.json" style="width: 400px; padding: 10px;" required>
				<br><br>
				<button type="submit" style="background: #28a745; color: white; padding: 10px 20px; border: none; border-radius: 5px;">Import Jobs</button>
			</form>
			<br>
			<p><small>💡 Export your local jobs first using: <code>go run ./cmd/export-jobs</code></small></p>
		</body></html>`))
		return
	}

	// Fetch JSON from URL
	resp, err := http.Get(jsonURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch URL: %v", err), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	// Parse JSON
	var jobs []JobImport
	err = json.NewDecoder(resp.Body).Decode(&jobs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Connect to database
	d, err := db.Connect()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Import jobs
	imported := 0
	for _, job := range jobs {
		err = d.InsertJobTyped(job.Title, job.Company, job.Location, job.Type, job.URL)
		if err == nil {
			imported++
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
	<!DOCTYPE html>
	<html><head><title>Jobs Imported</title></head>
	<body style="font-family: Arial; text-align: center; padding: 50px;">
		<h2>✅ Jobs Imported Successfully!</h2>
		<p><strong>%d out of %d jobs imported</strong></p>
		<p>(Duplicates are automatically skipped)</p>
		<a href="/login" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Go to Dashboard</a>
	</body></html>`, imported, len(jobs))))
}

// quickSetupHandler creates sample jobs for testing (for Render deployment)
func quickSetupHandler(w http.ResponseWriter, r *http.Request) {
	d, err := db.Connect()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database connection failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer d.Close()

	// Sample jobs to insert
	sampleJobs := []struct {
		title, company, location, jobType, url string
	}{
		{"Senior Software Engineer", "Y Combinator", "San Francisco, CA", "Full-time", "https://ycombinator.com/jobs/senior-swe"},
		{"Full Stack Developer", "OpenAI", "Remote", "Full-time", "https://openai.com/jobs/fullstack"},
		{"Backend Engineer", "Stripe", "San Francisco, CA", "Full-time", "https://stripe.com/jobs/backend"},
		{"Frontend Developer", "Airbnb", "San Francisco, CA", "Full-time", "https://airbnb.com/jobs/frontend"},
		{"DevOps Engineer", "Dropbox", "Remote", "Full-time", "https://dropbox.com/jobs/devops"},
		{"Data Scientist", "Uber", "San Francisco, CA", "Full-time", "https://uber.com/jobs/datascientist"},
		{"Product Manager", "Meta", "Menlo Park, CA", "Full-time", "https://meta.com/jobs/pm"},
		{"iOS Developer", "Apple", "Cupertino, CA", "Full-time", "https://apple.com/jobs/ios"},
		{"Machine Learning Engineer", "Google", "Mountain View, CA", "Full-time", "https://google.com/jobs/ml"},
		{"Security Engineer", "Netflix", "Los Gatos, CA", "Full-time", "https://netflix.com/jobs/security"},
	}

	// Insert sample jobs
	inserted := 0
	for _, job := range sampleJobs {
		err = d.InsertJobTyped(job.title, job.company, job.location, job.jobType, job.url)
		if err == nil {
			inserted++
		}
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(`
	<!DOCTYPE html>
	<html><head><title>Quick Setup Complete</title></head>
	<body style="font-family: Arial; text-align: center; padding: 50px;">
		<h2>✅ Quick Setup Complete!</h2>
		<p><strong>%d sample jobs added to your dashboard</strong></p>
		<p>You can now test the full functionality!</p>
		<a href="/login" style="background: #007bff; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Go to Dashboard</a>
	</body></html>`, inserted)))
}

// autoInitialize sets up admin user and sample jobs if database is empty
func autoInitialize(database *db.DB) {
	logger.Info("Starting database initialization check...")
	
	// Check if admin user exists
	_, _, _, err := database.GetUserByUsername("admin")
	if err != nil {
		logger.Info("Admin user not found, creating...")
		// Create admin user
		err = database.CreateUser("admin", "password123")
		if err != nil {
			logger.Error("Failed to create admin user: %v", err)
		} else {
			logger.Info("✅ Created admin user (admin/password123)")
		}
	} else {
		logger.Info("✅ Admin user already exists")
	}

	// Check if any jobs exist
	jobs, err := database.ListJobs(db.JobFilter{}, 1, 1)
	if err != nil {
		logger.Error("Failed to check existing jobs: %v", err)
		return
	}
	
	if len(jobs) == 0 {
		logger.Info("No jobs found, adding sample jobs...")
		// Add sample jobs
		sampleJobs := []struct {
			title, company, location, jobType, url string
		}{
			{"Senior Software Engineer", "Y Combinator", "San Francisco, CA", "Full-time", "https://ycombinator.com/jobs/senior-swe"},
			{"Full Stack Developer", "OpenAI", "Remote", "Full-time", "https://openai.com/jobs/fullstack"},
			{"Backend Engineer", "Stripe", "San Francisco, CA", "Full-time", "https://stripe.com/jobs/backend"},
			{"Frontend Developer", "Airbnb", "San Francisco, CA", "Full-time", "https://airbnb.com/jobs/frontend"},
			{"DevOps Engineer", "Dropbox", "Remote", "Full-time", "https://dropbox.com/jobs/devops"},
			{"Data Scientist", "Uber", "San Francisco, CA", "Full-time", "https://uber.com/jobs/datascientist"},
			{"Product Manager", "Meta", "Menlo Park, CA", "Full-time", "https://meta.com/jobs/pm"},
			{"iOS Developer", "Apple", "Cupertino, CA", "Full-time", "https://apple.com/jobs/ios"},
			{"Machine Learning Engineer", "Google", "Mountain View, CA", "Full-time", "https://google.com/jobs/ml"},
			{"Security Engineer", "Netflix", "Los Gatos, CA", "Full-time", "https://netflix.com/jobs/security"},
		}

		inserted := 0
		for _, job := range sampleJobs {
			err = database.InsertJobTyped(job.title, job.company, job.location, job.jobType, job.url)
			if err != nil {
				logger.Error("Failed to insert job %s: %v", job.title, err)
			} else {
				inserted++
			}
		}
		logger.Info("✅ Auto-initialized database with %d sample jobs", inserted)
	} else {
		logger.Info("✅ Database already contains %d jobs", len(jobs))
	}
	logger.Info("Database initialization complete")
}// -------------------- MAIN --------------------

func main() {
	// Check for PORT environment variable (required for Render)
	port := os.Getenv("PORT")
	if port == "" {
		// Fallback to flag or default
		portFlag := flag.String("port", "8080", "port to serve on")
		flag.Parse()
		port = *portFlag
	}

	// Connect to DB and ensure schema is created (from updated main.go logic)
	database, err := db.Connect()
	if err != nil {
		logger.Fatal("Failed to connect to DB: %v", err)
	}

	// Auto-initialize for Render deployment
	autoInitialize(database)

	// All routes are now based on the authenticated logic
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/logout", logoutHandler)

	// Admin initialization route (for Render deployment)
	http.HandleFunc("/init-admin", initAdminHandler)

	// Job import route (for Render deployment)
	http.HandleFunc("/import-jobs", importJobsHandler)

	// Quick setup route (for Render deployment)
	http.HandleFunc("/quick-setup", quickSetupHandler)

	// Protected routes (UI from main (1).go, logic from main.go)
	http.HandleFunc("/filters", AuthRequired(filtersHandler))
	http.HandleFunc("/dashboard", AuthRequired(dashboardHandler))
	http.HandleFunc("/results", AuthRequired(resultsHandler))
	http.HandleFunc("/download-csv", AuthRequired(downloadCSVHandler))
	http.HandleFunc("/mark-applied", AuthRequired(markAppliedHandler))

	logger.Info("Listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Server failed: %v", err)
	}
}
