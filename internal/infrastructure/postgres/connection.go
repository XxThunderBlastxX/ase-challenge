package postgres

import (
	"fmt"
	"log"

	"github.com/xxthunderblastxx/ase-challenge/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionManager struct {
	DB *gorm.DB
}

func MustConnect(cfg *config.PostgresConfig) *ConnectionManager {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("☹️ failed to connect database: %v", err)
	}

	log.Println("🚀 Connected to the database successfully")

	cm := &ConnectionManager{
		DB: db,
	}

	return cm
}
