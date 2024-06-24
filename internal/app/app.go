package app

import (
	"log/slog"

	"github.com/elusiv0/oz_task/pkg/httpserver"
)

type App struct {
	server *httpserver.HttpServer
	logger *slog.Logger
}

func New(
	server *httpserver.HttpServer,
	logger *slog.Logger,
) *App {
	app := &App{
		server: server,
		logger: logger,
	}

	return app
}

func (a *App) Run() error {
	a.logger.Info("starting http server...")
	if err := a.server.Start(); err != nil {
		a.logger.Error("error with up httpserver: " + err.Error())
		return err
	}
	return nil
}
