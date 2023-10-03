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
	http.ListenAndServe(":8080", nil)
}
