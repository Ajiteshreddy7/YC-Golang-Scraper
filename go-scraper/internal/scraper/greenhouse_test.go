package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScrapeGreenhouse(t *testing.T) {
	// Mock Greenhouse API response
	mockResp := `{
        "jobs": [
            {
                "title": "Software Engineer",
                "absolute_url": "https://example.com/jobs/123",
                "location": {"name": "San Francisco, CA"},
                "department": {"name": "Engineering"}
            },
            {
                "title": "Senior Staff Engineer",
                "absolute_url": "https://example.com/jobs/456",
                "location": {"name": "San Francisco, CA"},
                "department": {"name": "Engineering"}
            },
            {
                "title": "Software Engineer Intern",
                "absolute_url": "https://example.com/jobs/789",
                "location": {"name": "Remote, US"},
                "department": {"name": "Engineering"}
            }
        ]
    }`

	// Start test server that always returns our mock response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(mockResp))
	}))
	defer ts.Close()

	// Temporarily replace the Greenhouse API URL
	originalURL := greenhouseAPIURL
	greenhouseAPIURL = ts.URL + "/v1/boards/%s/jobs"
	defer func() { greenhouseAPIURL = originalURL }() // Run scraper
	jobs, err := ScrapeGreenhouse("test")
	if err != nil {
		t.Fatalf("ScrapeGreenhouse failed: %v", err)
	}

	// Should get 2 jobs (intern and regular engineer, not senior)
	if len(jobs) != 2 {
		t.Errorf("Expected 2 jobs, got %d", len(jobs))
	}

	// Check filtering worked
	for _, j := range jobs {
		if j.Title == "Senior Staff Engineer" {
			t.Errorf("Found senior role that should be filtered: %s", j.Title)
		}
	}
}

// mockTransport replaces API URL with test server URL
type mockTransport struct {
	origURL string
	mockURL string
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.String() == t.origURL {
		req.URL.Host = t.mockURL
		req.URL.Scheme = "http"
	}
	return http.DefaultTransport.RoundTrip(req)
}
