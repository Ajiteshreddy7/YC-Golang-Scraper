package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type greenhouseJob struct {
	Title    string `json:"title"`
	Absolute string `json:"absolute_url"`
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Department struct {
		Name string `json:"name"`
	} `json:"department"`
}

type greenhouseResponse struct {
	Jobs []greenhouseJob `json:"jobs"`
}

// Job represents the simplified job record used by DB layer
type Job struct {
	Title    string
	Company  string
	Location string
	URL      string
	Type     string
}

var seniorRe = regexp.MustCompile(`(?i)\b(senior|sr\.|lead|staff|principal|manager|director|architect|vp|head of|chief)\b`)
var earlyCareerRe = regexp.MustCompile(`(?i)\b(intern|internship|new grad|new graduate|associate|junior|entry level|entry-level|rotational|co-op|fellow|apprentice)\b`)
var usaLocs = []string{"united states", "usa", "us", "remote", "new york", "san francisco", "seattle", "austin", "boston", "chicago", "los angeles", "atlanta"}

func isEarlyCareer(title string) bool {
	t := strings.ToLower(title)
	if seniorRe.MatchString(t) {
		return false
	}
	if earlyCareerRe.MatchString(t) {
		return true
	}
	basicRoles := []string{"engineer", "developer", "analyst", "specialist", "coordinator"}
	for _, r := range basicRoles {
		if strings.Contains(t, r) {
			return true
		}
	}
	return false
}

func isInUSA(loc string) bool {
	l := strings.ToLower(loc)
	for _, id := range usaLocs {
		if strings.Contains(l, id) {
			return true
		}
	}
	return false
}

// API URL exposed for testing
var greenhouseAPIURL = "https://api.greenhouse.io/v1/boards/%s/jobs"

// ScrapeGreenhouse fetches and filters jobs for a given company identifier
func ScrapeGreenhouse(company string) ([]Job, error) {
	url := fmt.Sprintf(greenhouseAPIURL+"?content=true", company)
	client := &http.Client{Timeout: 20 * time.Second}

	var resp *http.Response
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		// Build request to set headers
		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("User-Agent", "yc-go-scraper/1.0 (+https://github.com/ajiteshreddy7/yc-go-scraper)")

		resp, err = client.Do(req)
		if err != nil {
			// network error: exponential backoff and retry
			time.Sleep(time.Duration(500*(1<<attempt)) * time.Millisecond)
			continue
		}

		// Handle rate limiting
		if resp.StatusCode == http.StatusTooManyRequests { // 429
			// Respect Retry-After if provided
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, convErr := strconv.Atoi(ra); convErr == nil {
					time.Sleep(time.Duration(secs) * time.Second)
				} else {
					time.Sleep(time.Duration(1<<attempt) * time.Second)
				}
			} else {
				time.Sleep(time.Duration(1<<attempt) * time.Second)
			}
			resp.Body.Close()
			continue
		}

		// Retry on 5xx
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			time.Sleep(time.Duration(500*(1<<attempt)) * time.Millisecond)
			continue
		}

		// Break for all other statuses (including 2xx and 4xx)
		break
	}

	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no response from greenhouse API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var gr greenhouseResponse
	if err := json.Unmarshal(body, &gr); err != nil {
		return nil, err
	}

	var out []Job
	for _, j := range gr.Jobs {
		title := j.Title
		loc := j.Location.Name
		if isEarlyCareer(title) && isInUSA(loc) {
			out = append(out, Job{
				Title:    title,
				Company:  strings.Title(company),
				Location: loc,
				URL:      j.Absolute,
				Type:     j.Department.Name,
			})
		}
	}
	return out, nil
}
