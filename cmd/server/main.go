package main

import (
	"chinese-game-backend/internal/config"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	router := chi.NewRouter()

	router.Get("/api/status", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	port := cfg.Port

	fmt.Printf("Start server")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
