package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	router.Get("/api/status", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	port := "8080"

	fmt.Printf("Start server")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
