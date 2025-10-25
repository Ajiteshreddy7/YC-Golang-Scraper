package exporter

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ajiteshreddy7/yc-go-scraper/internal/db"
)

// ExportCSV exports all rows from job_applications to the given CSV file
func ExportCSV(d *db.DB, path string) error {
	rows, err := d.Conn.Query(`SELECT title, company, location, type, url, date_added, status FROM job_applications ORDER BY date_added DESC`)
	if err != nil {
		return err
	}
	defer rows.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// header
	if err := w.Write([]string{"Title", "Company", "Location", "Type", "URL", "Date Added", "Status"}); err != nil {
		return err
	}

	for rows.Next() {
		var title, company, location, typ, url, dateAdded, status string
		if err := rows.Scan(&title, &company, &location, &typ, &url, &dateAdded, &status); err != nil {
			return err
		}
		if err := w.Write([]string{title, company, location, typ, url, dateAdded, status}); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	fmt.Printf("Wrote CSV to %s\n", path)
	return nil
}
