package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/response"
)

type ProductHandler struct {
	service product.Service
}

func NewProductHandler(s product.Service) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

func (h *ProductHandler) CreateProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p product.Product
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: err.Error(),
			})

		}

		if err := h.service.CreateProduct(&p); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
			Data: p,
		})
	}
}

func (h *ProductHandler) GetProductByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		p, err := h.service.GetProductByID(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorResponse{
				Error: "Product not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
			Data: p,
		})
	}
}

func (h *ProductHandler) GetAllProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		products, err := h.service.GetAllProducts()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
			Data: products,
		})
	}
}

func (h *ProductHandler) UpdateProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var p product.Product
		if err := c.BodyParser(&p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		if err := h.service.UpdateProduct(id, &p); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
			Data: p,
		})
	}
}

func (h *ProductHandler) DeleteProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if err := h.service.DeleteProduct(id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
				Error: err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
