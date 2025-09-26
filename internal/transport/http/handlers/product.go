package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
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

		// Parse request body
		if err := c.BodyParser(&p); err != nil {
			return errors.HandleError(c, errors.NewInvalidInputError("invalid request body format: "+err.Error()))
		}

		// Call service layer
		if err := h.service.CreateProduct(&p); err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleCreatedSuccess(c, p)
	}
}

func (h *ProductHandler) GetProductByID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate ID parameter
		if id == "" {
			return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
		}

		// Call service layer
		p, err := h.service.GetProductByID(id)
		if err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleSuccess(c, p)
	}
}

func (h *ProductHandler) GetAllProducts() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Call service layer
		products, err := h.service.GetAllProducts()
		if err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleSuccess(c, products)
	}
}

func (h *ProductHandler) UpdateProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate ID parameter
		if id == "" {
			return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
		}

		var p product.Product

		// Parse request body
		if err := c.BodyParser(&p); err != nil {
			return errors.HandleError(c, errors.NewInvalidInputError("invalid request body format: "+err.Error()))
		}

		// Call service layer
		if err := h.service.UpdateProduct(id, &p); err != nil {
			return errors.HandleError(c, err)
		}

		// Get updated product to return
		updatedProduct, err := h.service.GetProductByID(id)
		if err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleSuccess(c, updatedProduct)
	}
}

func (h *ProductHandler) DeleteProduct() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate ID parameter
		if id == "" {
			return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
		}

		// Call service layer
		if err := h.service.DeleteProduct(id); err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleNoContent(c)
	}
}

func (h *ProductHandler) IncrementStock() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate ID parameter
		if id == "" {
			return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
		}

		var req struct {
			StockIncrement int `json:"stock_increment" validate:"required,min=1"`
		}

		// Parse request body
		if err := c.BodyParser(&req); err != nil {
			return errors.HandleError(c, errors.NewInvalidInputError("invalid request body format: "+err.Error()))
		}

		// Validate increment value
		if req.StockIncrement <= 0 {
			return errors.HandleError(c, errors.NewInvalidInputError("stock_increment must be greater than 0"))
		}

		// Call service layer
		err := h.service.IncermentStock(id, req.StockIncrement)
		if err != nil {
			return errors.HandleError(c, err)
		}

		// Get updated product to return current stock
		updatedProduct, err := h.service.GetProductByID(id)
		if err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleSuccess(c, map[string]any{
			"message":          "Stock incremented successfully",
			"product":          updatedProduct,
			"increment_amount": req.StockIncrement,
		})
	}
}

func (h *ProductHandler) DecrementStock() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		// Validate ID parameter
		if id == "" {
			return errors.HandleError(c, errors.NewMissingRequiredDataError("id"))
		}

		var req struct {
			StockDecrement int `json:"stock_decrement" validate:"required,min=1"`
		}

		// Parse request body
		if err := c.BodyParser(&req); err != nil {
			return errors.HandleError(c, errors.NewInvalidInputError("invalid request body format: "+err.Error()))
		}

		// Validate decrement value
		if req.StockDecrement <= 0 {
			return errors.HandleError(c, errors.NewInvalidInputError("stock_decrement must be greater than 0"))
		}

		// Call service layer
		err := h.service.DecrementStock(id, req.StockDecrement)
		if err != nil {
			return errors.HandleError(c, err)
		}

		// Get updated product to return current stock
		updatedProduct, err := h.service.GetProductByID(id)
		if err != nil {
			return errors.HandleError(c, err)
		}

		return errors.HandleSuccess(c, map[string]any{
			"message":          "Stock decremented successfully",
			"product":          updatedProduct,
			"decrement_amount": req.StockDecrement,
		})
	}
}
