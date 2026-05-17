package main

import (
	"chinese-game-backend/internal/config"
	"chinese-game-backend/internal/repository"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	_, err := repository.NewDBConnection(cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)

	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	log.Println("БД подключена")

	router := chi.NewRouter()

	router.Get("/api/status", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})

	port := cfg.Port

	fmt.Printf("Запуск сервера на порте: %s", cfg.Port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
