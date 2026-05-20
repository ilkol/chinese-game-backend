package main

import (
	"chinese-game-backend/internal/app"
	"chinese-game-backend/internal/config"
	"log"
)

// @title Chinese Game API
// @description API for web game for learning Chinese
// @host localhost:5000
// @BasePAth /api
// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization
func main() {
	cfg := config.Load()

	app, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}

}
