package route

import (
	"api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App            *fiber.App
	UserController *http.UserController
	AuthMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// todo create login user
	c.App.Post("/api/login", c.UserController.Login)
	c.App.Post("/api/users", c.UserController.Register)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	//user
	c.App.Get("/api/users", c.UserController.FindAll)
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Put("/api/users/:id", c.UserController.Update)
	c.App.Delete("/api/users/:id", c.UserController.Delete)
}
