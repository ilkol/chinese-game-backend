package main

import (
	"chinese-game-backend/internal/app"
	"chinese-game-backend/internal/config"
	"log"
)

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
