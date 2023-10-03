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
