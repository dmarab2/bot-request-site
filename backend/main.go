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
	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("There was an error: %s\n", err.Error())
	}
}
