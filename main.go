package main

import (
    "net/http"
)

func main() {
    db, err := InitDB("./clicks.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    http.HandleFunc("/", RootHandler(db))

    // Specify the paths to your SSL certificate and private key files.
    certFile := "/path/to/your/ssl-certificate.pem"
    keyFile := "/path/to/your/private-key.pem"

    // Start an HTTPS server using your SSL certificate and private key.
    err = http.ListenAndServeTLS(":443", certFile, keyFile, nil)
    if err != nil {
        panic(err)
    }
}
