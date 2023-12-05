package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := InitDB("./clicks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", RootHandler(db))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
