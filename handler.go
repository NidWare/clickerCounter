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

		if r.URL.Path == "/favicon.ico" {
			http.NotFound(w, r)
			return
		}

		if r.Header.Get("Purpose") == "prefetch" {
			http.NotFound(w, r)
			return
		}

		referral := r.URL.Query().Get("utm_source")

		redirectURL := "https://fansly.com/FoxyFlair/posts"

		if referral != "" {
			redirects, err := GetRedirects(db)
			if err != nil {
				fmt.Println("Error caught while getting redirects")
			}
			if mappedURL, ok := redirects[referral]; ok {
				redirectURL = mappedURL
			}
		}

		cookie, err := r.Cookie("visited")

		if err != nil {
			if err == http.ErrNoCookie {
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
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if cookie.Value != "true" {
				fmt.Println("Strange cookies")
			}
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)

	}
}
