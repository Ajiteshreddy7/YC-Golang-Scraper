package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

// deriveLevels returns canonical level labels found in a job title
func deriveLevels(title string) []string {
	t := strings.ToLower(title)
	var out []string
	add := func(s string) { out = append(out, s) }
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
	if len(out) == 0 {
		if matched, _ := regexp.MatchString(`\b(engineer|developer|analyst|specialist|coordinator)\b`, t); matched {
			if ok, _ := regexp.MatchString(`\b(senior|staff|principal|lead|manager|director|architect|head|chief|vp)\b`, t); !ok {
				add("Entry Level")
			}
		}
	}
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

const indexHTML = `<!DOCTYPE html>
<html>
<head>
    <title>Job Opportunities Finder</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #f8f9fa; padding: 20px; }
        .container { max-width: 1200px; margin: 0 auto; }
        h1 { color: #343a40; text-align: center; margin-bottom: 10px; }
        .subtitle { text-align: center; color: #6c757d; margin-bottom: 30px; }
        .filters { background: white; padding: 24px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px; }
        .filter-section { margin-bottom: 20px; }
        .filter-label { font-weight: 600; color: #495057; margin-bottom: 8px; display: block; }
        .filter-row { display: flex; gap: 10px; flex-wrap: wrap; margin-bottom: 15px; }
        input[type="text"], select { padding: 10px 12px; border: 1px solid #ced4da; border-radius: 6px; flex: 1; min-width: 200px; font-size: 14px; }
        .level-checkboxes { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 10px; margin-top: 10px; }
        .level-checkboxes label { display: flex; align-items: center; gap: 6px; padding: 8px; background: #f8f9fa; border-radius: 4px; cursor: pointer; }
        .level-checkboxes input[type="checkbox"] { cursor: pointer; }
        .action-buttons { display: flex; gap: 10px; justify-content: center; margin-top: 20px; }
        button { background: #007bff; color: white; border: none; padding: 12px 24px; border-radius: 6px; cursor: pointer; font-weight: 500; font-size: 14px; }
        button:hover { background: #0069d9; }
        .btn-secondary { background: #6c757d; }
        .btn-secondary:hover { background: #5a6268; }
        .results-section { display: none; }
        .stats { display: flex; justify-content: space-around; margin-bottom: 20px; flex-wrap: wrap; gap: 15px; }
        .stat-card { background: white; padding: 16px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); text-align: center; min-width: 120px; }
        .stat-number { font-size: 1.8em; font-weight: bold; color: #007bff; }
        .stat-label { color: #6c757d; margin-top: 5px; font-size: 0.9em; }
        .jobs-list { background: white; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); overflow: hidden; }
        .job-item { padding: 16px; border-bottom: 1px solid #dee2e6; display: flex; justify-content: space-between; align-items: start; gap: 15px; }
        .job-item:hover { background: #f8f9fa; }
        .job-item:last-child { border-bottom: none; }
        .job-info { flex: 1; }
        .job-title { font-weight: 600; color: #343a40; margin-bottom: 6px; font-size: 1.05em; }
        .job-meta { color: #6c757d; font-size: 0.9em; line-height: 1.6; }
        .job-level { display: inline-block; background: #e7f3ff; color: #0056b3; padding: 2px 8px; border-radius: 12px; font-size: 0.85em; margin-top: 4px; }
        .job-actions { display: flex; flex-direction: column; gap: 8px; }
        .btn-apply { background: #28a745; padding: 8px 16px; border-radius: 6px; text-decoration: none; color: white; font-size: 0.9em; text-align: center; white-space: nowrap; }
        .btn-apply:hover { background: #218838; }
        .btn-mark-applied { background: #007bff; color: white; border: none; padding: 8px 16px; border-radius: 6px; cursor: pointer; font-size: 0.9em; white-space: nowrap; }
        .btn-mark-applied:hover { background: #0069d9; }
        .applied-badge { color: #28a745; font-weight: 600; font-size: 0.9em; white-space: nowrap; }
        .status-applied { opacity: 0.6; }
        .no-results { text-align: center; padding: 40px; color: #6c757d; }
        .select-all { margin-bottom: 10px; font-weight: normal; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🚀 Job Opportunities Finder</h1>
        <div class="subtitle">Find your next opportunity from {{.TotalJobs}} early-career positions</div>
        
        <div class="filters">
            <div class="filter-section">
                <label class="filter-label">Search</label>
                <input type="text" id="search" placeholder="Search by title, company, or location..." />
            </div>
            
            <div class="filter-section">
                <label class="filter-label">Job Levels</label>
                <div class="select-all">
                    <label><input type="checkbox" id="select-all" checked> Select / Deselect All</label>
                </div>
                <div class="level-checkboxes" id="levels">
                    {{range .Levels}}
                    <label><input type="checkbox" value="{{.}}" checked> {{.}}</label>
                    {{end}}
                </div>
            </div>
            
            <div class="filter-section">
                <label class="filter-label">Filters</label>
                <div class="filter-row">
                    <select id="company">
                        <option value="">All Companies</option>
                        {{range .Companies}}
                        <option value="{{.}}">{{.}}</option>
                        {{end}}
                    </select>
                    <select id="location">
                        <option value="">All Locations</option>
                        {{range .Locations}}
                        <option value="{{.}}">{{.}}</option>
                        {{end}}
                    </select>
                    <select id="status">
                        <option value="">All Statuses</option>
                        <option value="Not Applied">Not Applied</option>
                        <option value="Applied">Applied</option>
                    </select>
                </div>
            </div>
            
            <div class="action-buttons">
                <button onclick="filterJobs()">🔍 Show Jobs</button>
                <button class="btn-secondary" onclick="resetFilters()">↻ Reset Filters</button>
                <button class="btn-secondary" onclick="exportCSV()" style="display:none;" id="export-btn">⬇ Download CSV</button>
            </div>
        </div>
        
        <div class="results-section" id="results">
            <div class="stats">
                <div class="stat-card">
                    <div class="stat-number" id="filtered-count">0</div>
                    <div class="stat-label">Matching Jobs</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number" id="not-applied-count">0</div>
                    <div class="stat-label">Not Applied</div>
                </div>
                <div class="stat-card">
                    <div class="stat-number" id="applied-count">0</div>
                    <div class="stat-label">Applied</div>
                </div>
            </div>
            
            <div class="jobs-list" id="jobs-list">
                <div class="no-results">Click "Show Jobs" to see results</div>
            </div>
        </div>
    </div>
    
    <script>
        const allJobs = {{.JobsJSON}};
        
        // Load applied jobs from localStorage and update statuses
        (function initializeAppliedJobs() {
            const appliedJobs = JSON.parse(localStorage.getItem('appliedJobs') || '{}');
            allJobs.forEach(job => {
                if (appliedJobs[job.URL]) {
                    job.Status = 'Applied';
                }
            });
        })();
        
        document.getElementById('select-all').addEventListener('change', function() {
            const checkboxes = document.querySelectorAll('#levels input[type="checkbox"]');
            checkboxes.forEach(cb => cb.checked = this.checked);
        });
        
        function resetFilters() {
            document.getElementById('search').value = '';
            document.getElementById('company').value = '';
            document.getElementById('location').value = '';
            document.getElementById('status').value = '';
            document.getElementById('select-all').checked = true;
            document.querySelectorAll('#levels input[type="checkbox"]').forEach(cb => cb.checked = true);
            document.getElementById('results').style.display = 'none';
            document.getElementById('export-btn').style.display = 'none';
        }
        
        function filterJobs() {
            const search = document.getElementById('search').value.toLowerCase();
            const selectedLevels = Array.from(document.querySelectorAll('#levels input:checked')).map(cb => cb.value.toLowerCase());
            const company = document.getElementById('company').value;
            const location = document.getElementById('location').value;
            const status = document.getElementById('status').value;
            
            let filtered = allJobs.filter(job => {
                // Search filter
                if (search) {
                    const searchable = (job.Title + ' ' + job.Company + ' ' + job.Location).toLowerCase();
                    if (!searchable.includes(search)) return false;
                }
                
                // Level filter
                if (selectedLevels.length > 0) {
                    const jobLevels = job.Levels.toLowerCase();
                    const matchesLevel = selectedLevels.some(level => jobLevels.includes(level));
                    if (!matchesLevel) return false;
                }
                
                // Company filter
                if (company && job.Company !== company) return false;
                
                // Location filter
                if (location && job.Location !== location) return false;
                
                // Status filter
                if (status && job.Status !== status) return false;
                
                return true;
            });
            
            displayResults(filtered);
        }
        
        function displayResults(jobs) {
            const resultsSection = document.getElementById('results');
            const jobsList = document.getElementById('jobs-list');
            const exportBtn = document.getElementById('export-btn');
            
            resultsSection.style.display = 'block';
            exportBtn.style.display = jobs.length > 0 ? 'inline-block' : 'none';
            
            // Update stats
            updateStats(jobs);
            
            // Display jobs
            if (jobs.length === 0) {
                jobsList.innerHTML = '<div class="no-results">No jobs match your filters. Try adjusting your criteria.</div>';
                return;
            }
            
            jobsList.innerHTML = jobs.map((job, index) => 
                '<div class="job-item ' + (job.Status === 'Applied' ? 'status-applied' : '') + '" data-index="' + index + '">' +
                    '<div class="job-info">' +
                        '<div class="job-title">' + escapeHtml(job.Title) + '</div>' +
                        '<div class="job-meta">' +
                            '<strong>' + escapeHtml(job.Company) + '</strong> &bull; ' + escapeHtml(job.Location) + '<br>' +
                            'Added: ' + new Date(job.DateAdded).toLocaleDateString() + ' &bull; Status: <span id="status-' + index + '">' + job.Status + '</span>' +
                        '</div>' +
                        '<span class="job-level">' + escapeHtml(job.Levels) + '</span>' +
                    '</div>' +
                    '<div class="job-actions">' +
                        '<a href="' + escapeHtml(job.URL) + '" target="_blank" class="btn-apply">Apply &rarr;</a>' +
                        (job.Status === 'Not Applied' 
                            ? '<button class="btn-mark-applied" onclick="markApplied(' + index + ')">✓ Mark Applied</button>'
                            : '<span class="applied-badge">✓ Applied</span>') +
                    '</div>' +
                '</div>'
            ).join('');
            
            // Scroll to results
            resultsSection.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
        }
        
        function markApplied(index) {
            // Update the job status in memory
            const currentFiltered = getCurrentFilteredJobs();
            if (index < currentFiltered.length) {
                const job = currentFiltered[index];
                job.Status = 'Applied';
                
                // Store in localStorage
                const appliedJobs = JSON.parse(localStorage.getItem('appliedJobs') || '{}');
                appliedJobs[job.URL] = true;
                localStorage.setItem('appliedJobs', JSON.stringify(appliedJobs));
                
                // Update the display
                const statusSpan = document.getElementById('status-' + index);
                if (statusSpan) statusSpan.textContent = 'Applied';
                
                const jobItem = document.querySelector('[data-index="' + index + '"]');
                if (jobItem) {
                    jobItem.classList.add('status-applied');
                    const actionsDiv = jobItem.querySelector('.job-actions');
                    const markBtn = actionsDiv.querySelector('.btn-mark-applied');
                    if (markBtn) {
                        markBtn.outerHTML = '<span class="applied-badge">✓ Applied</span>';
                    }
                }
                
                // Update stats
                updateStats(currentFiltered);
            }
        }
        
        function getCurrentFilteredJobs() {
            const search = document.getElementById('search').value.toLowerCase();
            const selectedLevels = Array.from(document.querySelectorAll('#levels input:checked')).map(cb => cb.value.toLowerCase());
            const company = document.getElementById('company').value;
            const location = document.getElementById('location').value;
            const status = document.getElementById('status').value;
            
            return allJobs.filter(job => {
                if (search) {
                    const searchable = (job.Title + ' ' + job.Company + ' ' + job.Location).toLowerCase();
                    if (!searchable.includes(search)) return false;
                }
                if (selectedLevels.length > 0) {
                    const jobLevels = job.Levels.toLowerCase();
                    const matchesLevel = selectedLevels.some(level => jobLevels.includes(level));
                    if (!matchesLevel) return false;
                }
                if (company && job.Company !== company) return false;
                if (location && job.Location !== location) return false;
                if (status && job.Status !== status) return false;
                return true;
            });
        }
        
        function updateStats(jobs) {
            const notApplied = jobs.filter(j => j.Status === 'Not Applied').length;
            const applied = jobs.filter(j => j.Status === 'Applied').length;
            
            document.getElementById('filtered-count').textContent = jobs.length;
            document.getElementById('not-applied-count').textContent = notApplied;
            document.getElementById('applied-count').textContent = applied;
        }
        
        function escapeHtml(text) {
            const div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }
        
        function exportCSV() {
            const search = document.getElementById('search').value.toLowerCase();
            const selectedLevels = Array.from(document.querySelectorAll('#levels input:checked')).map(cb => cb.value.toLowerCase());
            const company = document.getElementById('company').value;
            const location = document.getElementById('location').value;
            const status = document.getElementById('status').value;
            
            let filtered = allJobs.filter(job => {
                if (search) {
                    const searchable = (job.Title + ' ' + job.Company + ' ' + job.Location).toLowerCase();
                    if (!searchable.includes(search)) return false;
                }
                if (selectedLevels.length > 0) {
                    const jobLevels = job.Levels.toLowerCase();
                    const matchesLevel = selectedLevels.some(level => jobLevels.includes(level));
                    if (!matchesLevel) return false;
                }
                if (company && job.Company !== company) return false;
                if (location && job.Location !== location) return false;
                if (status && job.Status !== status) return false;
                return true;
            });
            
            let csv = 'Date,Company,Title,Location,Level,Status,URL\n';
            filtered.forEach(job => {
                const row = [
                    new Date(job.DateAdded).toLocaleDateString(),
                    job.Company,
                    job.Title,
                    job.Location,
                    job.Levels,
                    job.Status,
                    job.URL
                ].map(field => '"' + String(field).replace(/"/g, '""') + '"');
                csv += row.join(',') + '\n';
            });
            
            const blob = new Blob([csv], { type: 'text/csv' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'filtered_jobs.csv';
            a.click();
            window.URL.revokeObjectURL(url);
        }
        
        // Allow Enter key to trigger search
        document.getElementById('search').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') filterJobs();
        });
    </script>
</body>
</html>`

func main() {
	outDir := flag.String("out", "public", "Output directory for static site")
	flag.Parse()

	logger.InitFromEnv()
	logger.Info("Generating static site")

	d, err := db.Connect()
	if err != nil {
		logger.Fatal("db connect: %v", err)
	}
	defer d.Close()

	// Fetch all jobs
	rows, err := d.Conn.Query(`SELECT id, title, company, location, type, url, date_added, status FROM job_applications ORDER BY date_added DESC`)
	if err != nil {
		logger.Fatal("query jobs: %v", err)
	}
	defer rows.Close()

	type JobWithLevels struct {
		Job
		Levels      string
		StatusClass string
	}

	var jobs []JobWithLevels
	levelSet := map[string]bool{}
	companySet := map[string]bool{}
	locationSet := map[string]bool{}
	notApplied := 0
	applied := 0

	for rows.Next() {
		var job Job
		var typ string
		if err := rows.Scan(&job.ID, &job.Title, &job.Company, &job.Location, &typ, &job.URL, &job.DateAdded, &job.Status); err != nil {
			logger.Error("scan row: %v", err)
			continue
		}
		job.Type = typ

		// Derive levels
		levels := deriveLevels(job.Title)
		for _, lv := range levels {
			levelSet[lv] = true
		}
		levelsStr := strings.Join(levels, ", ")
		if levelsStr == "" {
			levelsStr = "General"
		}

		statusClass := "not-applied"
		if job.Status == "Applied" {
			statusClass = "applied"
			applied++
		} else {
			notApplied++
		}

		jobs = append(jobs, JobWithLevels{
			Job:         job,
			Levels:      levelsStr,
			StatusClass: statusClass,
		})

		companySet[job.Company] = true
		locationSet[job.Location] = true
	}

	// Convert sets to sorted slices
	var levels []string
	for k := range levelSet {
		levels = append(levels, k)
	}
	sort.Strings(levels)

	var companies []string
	for k := range companySet {
		companies = append(companies, k)
	}
	sort.Strings(companies)

	var locations []string
	for k := range locationSet {
		locations = append(locations, k)
	}
	sort.Strings(locations)

	// Generate index.html
	tmpl := template.Must(template.New("index").Parse(indexHTML))

	// Convert jobs to JSON for client-side filtering
	jobsJSON, err := json.Marshal(jobs)
	if err != nil {
		logger.Fatal("marshal jobs json: %v", err)
	}

	data := struct {
		Jobs       []JobWithLevels
		Levels     []string
		Companies  []string
		Locations  []string
		TotalJobs  int
		NotApplied int
		Applied    int
		JobsJSON   template.JS
	}{
		Jobs:       jobs,
		Levels:     levels,
		Companies:  companies,
		Locations:  locations,
		TotalJobs:  len(jobs),
		NotApplied: notApplied,
		Applied:    applied,
		JobsJSON:   template.JS(jobsJSON),
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		logger.Fatal("create output dir: %v", err)
	}

	indexPath := filepath.Join(*outDir, "index.html")
	f, err := os.Create(indexPath)
	if err != nil {
		logger.Fatal("create index.html: %v", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		logger.Fatal("execute template: %v", err)
	}

	logger.Info("Generated static site in %s", *outDir)

	// Also export jobs.json for API access
	jobsJSONPath := filepath.Join(*outDir, "jobs.json")
	jf, err := os.Create(jobsJSONPath)
	if err != nil {
		logger.Fatal("create jobs.json: %v", err)
	}
	defer jf.Close()

	enc := json.NewEncoder(jf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		logger.Fatal("encode json: %v", err)
	}

	logger.Info("Generated jobs.json")
	logger.Info("Site ready with %d jobs", len(jobs))
}
