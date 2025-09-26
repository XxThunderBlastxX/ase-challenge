package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/config"
	"github.com/xxthunderblastxx/ase-challenge/internal/infrastructure/postgres"
)

// App struct holds the web-server configuration
type App struct {
	*fiber.App

	Appconfig    *config.AppConfig
	PostgresConn *postgres.ConnectionManager
}

// New initializes and returns a new Server instance
func New() *App {
	cfg := config.New()

	pConn := postgres.MustConnect(&cfg.PostgresConfig)

	return &App{
		App:          fiber.New(),
		Appconfig:    config.New(),
		PostgresConn: pConn,
	}
}
