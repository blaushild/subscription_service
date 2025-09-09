package main

import (
	"log"
	"subscribe_service/internal/config"
	"subscribe_service/run"
)

// @title Subscription service API
// @version 0.0.1
// @description This is a subscription service.
// @BasePath /api/v1/subscription
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app := run.NewApp(cfg)
	app.Run()
}
