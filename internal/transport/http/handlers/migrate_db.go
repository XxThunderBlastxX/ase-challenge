package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	"github.com/xxthunderblastxx/ase-challenge/internal/infrastructure/postgres"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/response"
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
		if err := h.conn.DB.AutoMigrate(product.Product{}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
			Data: "Database migrated successfully!!",
		})
	}
}
