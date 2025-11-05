package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/exporter"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/logger"
	"github.com/ajiteshreddy7/yc-go-scraper/internal/scraper"
)

type Config struct {
	TargetPlatforms map[string][]string `json:"target_platforms"`
}

func main() {
	// CLI flags
	cfgPath := flag.String("config", "config/scraper_config.json", "Path to scraper config JSON")
	outPath := flag.String("out", "data/job_applications.csv", "Path to output CSV file")
	flag.Parse()

	// Init logger level from env
	logger.InitFromEnv()

	logger.Info("Starting Go Job Scraper")

	// Connect to DB
	d, err := db.Connect()
	if err != nil {
		logger.Fatal("db connect: %v", err)
	}
	defer d.Close()

	if err := d.CreateSchema(); err != nil {
		logger.Fatal("create schema: %v", err)
	}

	// Load config
	if _, err := os.Stat(*cfgPath); os.IsNotExist(err) {
		logger.Fatal("config file not found: %s", *cfgPath)
	}

	raw, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		logger.Fatal("read config: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(raw, &cfg); err != nil {
		logger.Fatal("unmarshal config: %v", err)
	}

	total := 0

	// Process Greenhouse if configured
	if companies, ok := cfg.TargetPlatforms["greenhouse"]; ok {
		logger.Info("Found %d greenhouse companies to scrape", len(companies))
		for i, c := range companies {
			logger.Info("[%d/%d] scraping %s (Greenhouse)", i+1, len(companies), c)
			jobs, err := scraper.ScrapeGreenhouse(c)
			if err != nil {
				logger.Warn("error scraping %s: %v", c, err)
				continue
			}
			for _, job := range jobs {
				if err := d.InsertJobTyped(job.Title, job.Company, job.Location, job.Type, job.URL); err != nil {
					logger.Error("insert job error: %v", err)
				} else {
					total++
				}
			}
			// be respectful
			time.Sleep(2 * time.Second)
		}
		logger.Info("Processed %d greenhouse jobs", total)
	}

	// Process Lever if configured
	if companies, ok := cfg.TargetPlatforms["lever"]; ok {
		logger.Info("Found %d lever companies to scrape", len(companies))
		leverTotal := 0
		for i, c := range companies {
			logger.Info("[%d/%d] scraping %s (Lever)", i+1, len(companies), c)
			jobs, err := scraper.ScrapeLever(c)
			if err != nil {
				logger.Warn("error scraping %s: %v", c, err)
				continue
			}
			for _, job := range jobs {
				if err := d.InsertJobTyped(job.Title, job.Company, job.Location, job.Type, job.URL); err != nil {
					logger.Error("insert job error: %v", err)
				} else {
					leverTotal++
					total++
				}
			}
			// be respectful
			time.Sleep(2 * time.Second)
		}
		logger.Info("Processed %d lever jobs", leverTotal)
	}
	
	logger.Info("Processed %d total jobs", total)

	// Export CSV
	if err := os.MkdirAll(filepath.Dir(*outPath), 0755); err != nil {
		logger.Fatal("create output dir: %v", err)
	}

	if err := exporter.ExportCSV(d, *outPath); err != nil {
		logger.Fatal("export csv: %v", err)
	}
	logger.Info("Exported CSV to %s", *outPath)
}
