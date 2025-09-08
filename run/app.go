package run

import (
	"subscribe_service/internal/config"
	"subscribe_service/internal/controller"
	"subscribe_service/internal/server"
)

type app struct {
	cfg        *config.Config
	httpServer *server.Server
}

func NewApp(cfg *config.Config) *app {
	c := controller.NewController(cfg)
	return &app{
		cfg:        cfg,
		httpServer: server.NewServer(cfg, c),
	}
}

func (a *app) Run() {
	a.httpServer.Serve()
}
