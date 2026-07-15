package app

import (
	levelv1 "chinese-game-backend/api/gen/level"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "chinese-game-backend/docs"
)

type App struct {
	Config          *config.Config
	DB              *sqlx.DB
	GrpcLevelClient levelv1.LevelServiceClient
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

	GrpcLevelClient, err := newGrpcLevelClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Ошибка инициализации gRPC клиента %v", err)
	}
	log.Printf("gRPC клиент подключен")

	return &App{cfg, db, GrpcLevelClient}, nil
}

func newGrpcLevelClient(cfg *config.Config) (levelv1.LevelServiceClient, error) {
	grpcAddr := cfg.LEVEL_SERVICE_ADDRES

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	con, err := grpc.NewClient(
		grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("Не удалось подключиться к gRPC серверу: %v", err)
	}

	con.Connect()

	for {
		state := con.GetState()
		if state.String() == "READY" {
			break // Успешно подключились!
		}

		// На каждом шаге цикла проверяем, не истекли ли наши 5 секунд
		if !con.WaitForStateChange(ctx, state) {
			return nil, fmt.Errorf("таймаут подключения к gRPC серверу: %w", ctx.Err())
		}
	}

	client := levelv1.NewLevelServiceClient(con)
	return client, nil
}

func (app *App) Run() error {
	userRepo := repository.NewUserRepository(app.DB)
	levelRepo := repository.NewLevelRepository(app.DB)
	progressRepo := repository.NewProgressRepository(app.DB)

	userService := service.NewUserService(userRepo, app.Config.JWTSecret)
	levelService := service.NewLevelService(levelRepo)
	progressService := service.NewProgressService(progressRepo)

	authHandler := handlers.NewAuthHandler(userService)
	levelHandler := handlers.NewLevelHandler(levelService)
	progressHandler := handlers.NewProgressHandler(progressService)
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

			r.Post("/progress", progressHandler.CompleteStep)
			r.Get("/progress/levels", progressHandler.GetCompletedLevels)
			r.Get("/progress/levels/{level_id}", progressHandler.IsLevelCompleted)

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
