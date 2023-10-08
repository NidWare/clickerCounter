package main

import (
	"context"
	"log"
	"net/http"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := InitDB("./clicks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", RootHandler(db))

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatal(err)
		}
	}()

	go StartTelegramBot(ctx, db)

	certFile := "./fullchain.pem"
	keyFile := "./privkey.pem"

	if err := http.ListenAndServeTLS(":443", certFile, keyFile, nil); err != nil {
		log.Fatal(err)
	}
}
