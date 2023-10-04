package main

import (
	"net/http"
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
}

func main() {
	db, err := InitDB("./clicks.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", RootHandler(db))

	// Redirect HTTP to HTTPS.
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(redirectTLS)); err != nil {
			panic(err)
		}
	}()

	// Specify the paths to your SSL certificate and private key files.
	certFile := "./ssl-certificate.pem"
	keyFile := "./private-key.pem"

	// Start an HTTPS server using your SSL certificate and private key.
	if err := http.ListenAndServeTLS(":443", certFile, keyFile, nil); err != nil {
		panic(err)
	}
}
