package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	"github.com/xxthunderblastxx/ase-challenge/internal/infrastructure/postgres"
	"github.com/xxthunderblastxx/ase-challenge/internal/transport/http/handlers"
)

func (r *Router) productRouter(grp fiber.Router) {
	pgrp := grp.Group("/products")

	repo := postgres.NewProductRepository(r.app.PostgresConn)
	s := product.NewService(repo)
	h := handlers.NewProductHandler(s)

	{
		pgrp.Post("/", h.CreateProduct())
		pgrp.Get("/", h.GetAllProducts())
		pgrp.Get("/:id", h.GetProductByID())
		pgrp.Put("/:id", h.UpdateProduct())
		pgrp.Delete("/:id", h.DeleteProduct())

		pgrp.Post("/:id/increment-stock", h.IncrementStock())
		pgrp.Post("/:id/decrement-stock", h.DecrementStock())
	}
}
