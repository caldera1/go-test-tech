package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	fmt.Println("server running on :8080")
	_ = http.ListenAndServe(":8080", nil)
}
