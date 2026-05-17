package main

import (
	"chinese-game-backend/internal/config"
	"chinese-game-backend/internal/repository"
	"chinese-game-backend/internal/service"
	handlers "chinese-game-backend/internal/transport/http"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)
	if err := repository.RunMigrations(dbURL); err != nil {
		log.Fatalf("Миграции БД провалены: %v", err)
	}

	db, err := repository.NewDBConnection(cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)

	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	log.Println("БД подключена")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(userService)

	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})
		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("ok"))
		})
	})

	port := cfg.Port

	fmt.Printf("Запуск сервера на порте: %s\n", cfg.Port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("%s\n", err)
	}
}
