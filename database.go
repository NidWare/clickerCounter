package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS clicks (referral TEXT PRIMARY KEY, clicks INTEGER)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func UpdateClicks(db *sql.DB, referral string) error {
	_, err := db.Exec(`INSERT OR IGNORE INTO clicks(referral, clicks) VALUES(?, 1); UPDATE clicks SET clicks = clicks + 1 WHERE referral = ?`, referral, referral)
	return err
}

func GetRedirects(db *sql.DB) (map[string]string, error) {
	redirects := make(map[string]string)

	rows, err := db.Query("SELECT param, value FROM redirects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through rows
	for rows.Next() {
		var key, url string
		if err := rows.Scan(&key, &url); err != nil {
			return nil, err
		}
		redirects[key] = url
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return redirects, nil
}
