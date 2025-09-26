package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/watchakorn-18k/scalar-go"
	"github.com/xxthunderblastxx/ase-challenge/internal/server"
	"github.com/xxthunderblastxx/ase-challenge/internal/transport/http/handlers"
)

type Router struct {
	app *server.App
}

func NewRouter(s *server.App) *Router {
	return &Router{
		app: s,
	}
}

func (r *Router) RegisterRoutes() {
	// Middleware
	r.app.Use(cors.New())
	r.app.Use(recover.New())

	// API Documentation route
	r.app.Use("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Product Inventory Management API Docs",
			},
			DarkMode: true,
		})

		if err != nil {
			return err
		}
		c.Type("html")
		return c.SendString(htmlContent)
	})

	// Base Group
	g := r.app.Group("/api/v1")

	// Register other routes here
	r.migrateDBRouter(g)
	r.productRouter(g)
}

func (r *Router) migrateDBRouter(grp fiber.Router) {
	dbConn := r.app.PostgresConn

	h := handlers.NewMigrateDBHandler(dbConn)

	grp.Get("/migrate", h.MigrateDB())
}
