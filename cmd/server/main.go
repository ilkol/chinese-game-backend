package main

import (
	"chinese-game-backend/internal/config"
	"chinese-game-backend/internal/repository"
	"chinese-game-backend/internal/service"
	handlers "chinese-game-backend/internal/transport/http"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	levelRepo := repository.NewLevelRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	levelService := service.NewLevelService(levelRepo)
	progressService := service.NewProgressRepository(progressRepo)

	authHandler := handlers.NewAuthHandler(userService)
	levelHandler := handlers.NewLevelHandler(levelService, progressService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("ok"))
		})

		r.Group(func(r chi.Router) {
			r.Use(authHandler.UserIdentity)

			r.Get("/level", levelHandler.GetAll)
			r.Get("/level/{id}", levelHandler.GetByID)

			r.Post("/progress", levelHandler.CompleteStep)
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Запуск сервера на порте: %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("%s\n", err)
		}
	}()

	<-quit
	log.Println("Остановка сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Сервер завершает работу с ошибкой: %v", err)
	}

	log.Println("Закрытие соединение с БД...")
	if err := db.Close(); err != nil {
		log.Printf("Ошибка закрытия соединения: %v", err)
	}

	log.Println("Сервер завершил работу")
}
