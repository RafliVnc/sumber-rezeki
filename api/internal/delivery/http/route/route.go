package route

import (
	"api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
}

func (c *RouteConfig) Setup() {
	//user
	c.App.Get("/api/users", c.UserController.FindAll)
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Put("/api/users/:id", c.UserController.Update)
	c.App.Delete("/api/users/:id", c.UserController.Delete)
}
