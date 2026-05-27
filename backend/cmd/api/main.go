package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"status": "ok", "message": "CI/CD loop closed. Backend is running."}`)
	})
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}