package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type leverJob struct {
	Text     string `json:"text"`      // Job title
	Hostedurl string `json:"hostedUrl"` // Job posting URL
	Categories struct {
		Location     string `json:"location"`
		Commitment   string `json:"commitment"`
		Team        string `json:"team"`
		Level       string `json:"level"`
	} `json:"categories"`
}

// API URL for Lever
var leverAPIURL = "https://api.lever.co/v0/postings/%s?mode=json"

// ScrapeLever fetches and filters jobs from Lever's API
func ScrapeLever(company string) ([]Job, error) {
	url := fmt.Sprintf(leverAPIURL, company)
	client := &http.Client{Timeout: 20 * time.Second}

	var resp *http.Response
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		req, reqErr := http.NewRequest("GET", url, nil)
		if reqErr != nil {
			return nil, reqErr
		}
		req.Header.Set("User-Agent", "yc-go-scraper/1.0 (+https://github.com/ajiteshreddy7/yc-go-scraper)")

		resp, err = client.Do(req)
		if err != nil {
			time.Sleep(time.Duration(500*(1<<attempt)) * time.Millisecond)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			time.Sleep(time.Duration(1<<attempt) * time.Second)
			continue
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			time.Sleep(time.Duration(500*(1<<attempt)) * time.Millisecond)
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("no response from lever API")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var jobs []leverJob
	if err := json.Unmarshal(body, &jobs); err != nil {
		return nil, err
	}

	var out []Job
	for _, j := range jobs {
		title := j.Text
		loc := j.Categories.Location
		if isEarlyCareer(title) && isInUSA(loc) {
			out = append(out, Job{
				Title:    title,
				Company:  strings.Title(company),
				Location: loc,
				URL:      j.Hostedurl,
				Type:     j.Categories.Team,
			})
		}
	}
	return out, nil
}