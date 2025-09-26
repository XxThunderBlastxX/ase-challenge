package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/config"
)

// Server struct holds the web-server configuration
type Server struct {
	*fiber.App

	Appconfig *config.AppConfig
}

// New initializes and returns a new Server instance
func New() *Server {
	return &Server{
		App:       fiber.New(),
		Appconfig: config.New(),
	}
}
