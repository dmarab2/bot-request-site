package main

import (
	"fmt"
	"net/http"
)

// A simple beginning function for now.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you've requested %s\n", r.URL.Path)
	})
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Printf("There was an error: %s\n", err.Error())
	}
}
