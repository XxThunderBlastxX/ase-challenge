package main

import (
	"log"

	"github.com/xxthunderblastxx/ase-challenge/internal/server"
)

func main() {
	s := server.New()

	// Run the application
	if err := s.Listen(":" + s.Appconfig.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
