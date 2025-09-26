package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// CORS middleware
	r.app.Use(cors.New())

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
