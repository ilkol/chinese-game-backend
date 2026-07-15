package main

import (
	levelv1 "chinese-game-backend/api/gen/level"
	"chinese-game-backend/internal/config"
	"chinese-game-backend/internal/repository"
	"chinese-game-backend/internal/service"
	"chinese-game-backend/internal/transport/grpc/level"
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewDBConnection(cfg.DB_User, cfg.DB_Pass, cfg.DB_Host, cfg.DB_Port, cfg.DB_Name)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	levelRepo := repository.NewLevelRepository(db)
	levelService := service.NewLevelService(levelRepo)

	grpcServer := grpc.NewServer()

	levelServer := level.NewLevelServer(levelService)
	levelv1.RegisterLevelServiceServer(grpcServer, levelServer)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка создания слушателя %v", err)
	}

	log.Printf("gRPC сервис уровней запущен на порту :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка запуска gRPC сервера %v", err)
	}
}
