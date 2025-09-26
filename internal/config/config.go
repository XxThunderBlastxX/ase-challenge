package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// PostgresConfig holds the configuration details for connecting to a PostgreSQL database.
type PostgresConfig struct {
	Host     string
	User     string
	Password string
	DB       string
	Port     string
}

// AppConfig holds the application wide configuration.
type AppConfig struct {
	Port   string
	Uptime time.Time

	PostgresConfig PostgresConfig
}

// New reads the .env file and returns an AppConfig instance populated with environment variables.
func New() *AppConfig {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file : " + err.Error())
	}

	return &AppConfig{
		Port:   os.Getenv("PORT"),
		Uptime: time.Now(),
		PostgresConfig: PostgresConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
			Port:     os.Getenv("POSTGRES_PORT"),
		},
	}
}
