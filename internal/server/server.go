package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/config"
	"github.com/xxthunderblastxx/ase-challenge/internal/infrastructure/postgres"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
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
		App: fiber.New(fiber.Config{
			ErrorHandler: errors.ErrorHandler(),
		}),
		Appconfig:    config.New(),
		PostgresConn: pConn,
	}
}
