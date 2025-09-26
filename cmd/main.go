package main

import (
	"log"

	"github.com/xxthunderblastxx/ase-challenge/internal/server"
	"github.com/xxthunderblastxx/ase-challenge/internal/transport/http/router"
)

func main() {
	// Initialize server
	s := server.New()

	// Register routes
	router.NewRouter(s).RegisterRoutes()

	// Run the application
	if err := s.Listen(":" + s.Appconfig.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
