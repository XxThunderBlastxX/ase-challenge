package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/config"
)

// App struct holds the web-server configuration
type App struct {
	*fiber.App

	Appconfig *config.AppConfig
}

// New initializes and returns a new Server instance
func New() *App {
	return &App{
		App:       fiber.New(),
		Appconfig: config.New(),
	}
}
