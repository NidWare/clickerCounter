package main

import (
	"database/sql"
	"net/http"
	"time"
)

// Define a map of redirects based on the "r" parameter
var redirects = map[string]string{
	"kj_bennet": "https://fans.ly/subscriptions/giftcode/NTY3MDk5Nzg3NDM4OTkzNDEwOjE6MTowYjMwNmY1NzA3",
	// Add more mappings as needed
}

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
			// Check if the "r" parameter is in the redirects map
			if redirectURL, ok := redirects[referral]; ok {
				http.Redirect(w, r, redirectURL, http.StatusSeeOther)
			} else {
				// Redirect to a default URL if "r" is not found in the map
				http.Redirect(w, r, "https://fansly.com/VikiMinelli/posts", http.StatusSeeOther)
			}

			expire := time.Now().AddDate(1, 0, 0)
			http.SetCookie(w, &http.Cookie{
				Name:    "visited",
				Value:   "true",
				Expires: expire,
			})
		}
	}
}
