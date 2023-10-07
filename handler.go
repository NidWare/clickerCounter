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

		// Always redirect to the specified URL
		redirectURL := "https://google.com" // Change this to your default redirect URL
		if referral != "" {
			// Define a map of redirects based on the "r" parameter
			redirects := map[string]string{
				"king":    "https://king.com",
				"example": "https://example.com",
				// Add more mappings as needed
			}

			// Check if the "r" parameter is in the redirects map
			if mappedURL, ok := redirects[referral]; ok {
				redirectURL = mappedURL
			}
		}

		// Redirect to the determined URL
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)

		if err != nil || cookie.Value != "true" {
			// Update clicks only if the user hasn't visited before
			err = UpdateClicks(db, referral)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Set the visited cookie to indicate that the user has visited
			expire := time.Now().AddDate(1, 0, 0)
			http.SetCookie(w, &http.Cookie{
				Name:    "visited",
				Value:   "true",
				Expires: expire,
			})
		}
	}
}
