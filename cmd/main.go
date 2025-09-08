package main

import (
	"log"
	"subscribe_service/internal/config"
	"subscribe_service/run"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app := run.NewApp(cfg)
	app.Run()
}
