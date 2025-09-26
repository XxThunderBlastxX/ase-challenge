package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	"github.com/xxthunderblastxx/ase-challenge/internal/infrastructure/postgres"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
)

type MigrateDBHandler struct {
	conn *postgres.ConnectionManager
}

func NewMigrateDBHandler(conn *postgres.ConnectionManager) *MigrateDBHandler {
	return &MigrateDBHandler{
		conn: conn,
	}
}

func (h *MigrateDBHandler) MigrateDB() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate connection manager
		if h.conn == nil {
			return errors.HandleError(c, errors.NewConnectionError("database connection manager is not initialized"))
		}

		// Validate database connection
		if h.conn.DB == nil {
			return errors.HandleError(c, errors.NewConnectionError("database connection is not established"))
		}

		// Perform migration
		if err := h.conn.DB.AutoMigrate(product.Product{}); err != nil {
			return errors.HandleError(c, errors.NewMigrationError("failed to migrate database: "+err.Error()))
		}

		return errors.HandleSuccess(c, "Database migrated successfully", fiber.StatusOK)
	}
}
