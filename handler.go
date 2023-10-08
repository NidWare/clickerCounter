package main

import (
	"database/sql"
	"fmt"
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

		redirectURL := "https://fansly.com/VikiMinelli/posts" // Change this to your default redirect URL
		if referral != "" {
			redirects, err := GetRedirects(db)

			if err != nil {
				fmt.Println("Error caught while getting redirects")
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
