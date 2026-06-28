package app

import (
	"chinese-game-backend/internal/config"
	"chinese-game-backend/internal/domain"
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
	"github.com/jmoiron/sqlx"
	"github.com/rs/cors"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "chinese-game-backend/docs"
)

type App struct {
	Config *config.Config
	DB     *sqlx.DB
}

func NewApp(cfg *config.Config) (*App, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)
	if err := repository.RunMigrations(dbURL); err != nil {
		return nil, fmt.Errorf("Миграции БД провалены: %w", err)
	}

	db, err := repository.NewDBConnection(cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)

	if err != nil {
		return nil, fmt.Errorf("Ошибка инициализации БД: %w", err)
	}
	log.Println("БД подключена")

	return &App{cfg, db}, nil
}

func (app *App) Run() error {
	userRepo := repository.NewUserRepository(app.DB)
	levelRepo := repository.NewLevelRepository(app.DB)
	progressRepo := repository.NewProgressRepository(app.DB)

	userService := service.NewUserService(userRepo, app.Config.JWTSecret)
	levelService := service.NewLevelService(levelRepo)
	progressService := service.NewProgressRepository(progressRepo)

	authHandler := handlers.NewAuthHandler(userService)
	levelHandler := handlers.NewLevelHandler(levelService, progressService)
	userHandler := handlers.NewUserHandler(userService)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://z9128573.beget.tech"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler)

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

			r.Group(func(r chi.Router) {
				r.Use(authHandler.CheckRole(domain.RoleStudent))

				r.Post("/user/join", userHandler.JoinStudentToTeacher)
			})

			r.Group(func(r chi.Router) {
				r.Use(authHandler.CheckRole(domain.RoleTeacher))

				r.Get("/teacher/students", userHandler.GetMyStudents)
				r.Get("/teacher/invite-code", userHandler.GetMyInviteCode)
			})

			r.Group(func(r chi.Router) {
				r.Use(authHandler.CheckRole(domain.RoleAdmin))

				r.Post("/level/{id}/step", levelHandler.CreateStep)
				r.Put("/level/{id}/step/{step_id}", levelHandler.UpdateStep)
				r.Delete("/level/{id}/step/{step_id}", levelHandler.DeleteStep)
				r.Put("/step/{step_id}/dialog", levelHandler.UpsertDialog)
			})
		})
	})

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	return app.startHTPPServer(router)
}

func (app *App) startHTPPServer(router *chi.Mux) error {
	defer func() {
		log.Println("Закрытие соединение с БД...")
		if err := app.DB.Close(); err != nil {
			log.Printf("%v\n", err)
		}
	}()

	srv := &http.Server{
		Addr:    ":" + app.Config.Port,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	errCh := make(chan error)

	go func() {
		log.Printf("Запуск сервера на порте: %s\n", app.Config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-quit:
		log.Println("Получен сигнал остановки")
	}
	log.Println("Остановка сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("Ошибка остановки сервера: %w", err)
	}

	fmt.Println("Сервер завершил работу")
	return nil
}
