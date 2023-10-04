package main

import (
	"database/sql"
	"net/http"
	"time"
)

func RootHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userAgent := r.Header.Get("User-Agent")
		if CheckIfBot(userAgent) {
			http.Error(w, "Bot detected", http.StatusForbidden)
			return
		}

		referral := r.URL.Query().Get("r")
		cookie, err := r.Cookie("visited")

		if err != nil || cookie.Value != "true" {
			err = UpdateClicks(db, referral)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			expire := time.Now().AddDate(1, 0, 0)
			http.SetCookie(w, &http.Cookie{
				Name:    "visited",
				Value:   "true",
				Expires: expire,
			})
		}

		// Redirect to Google.com
		http.Redirect(w, r, "google.com", http.StatusSeeOther)
	}
}
