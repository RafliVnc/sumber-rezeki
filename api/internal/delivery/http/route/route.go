package route

import (
	"api/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                          *fiber.App
	UserController               *http.UserController
	SalesController              *http.SalesController
	HelloController              *http.HelloController
	RouteController              *http.RouteController
	EmployeeController           *http.EmployeeController
	EmployeeAttendanceController *http.EmployeeAttendanceController
	AuthMiddleware               fiber.Handler
	Config                       *viper.Viper
}

func (c *RouteConfig) Setup() {
	// set cors
	c.App.Use(cors.New(cors.Config{
		AllowOrigins: c.Config.GetString("CORS_ALLOW_ORIGINS"),
	}))

	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// login user
	c.App.Post("/api/login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	//user
	c.App.Get("/api/hello", c.HelloController.SayHello)
	c.App.Get("/api/current", c.UserController.Current)
	c.App.Get("/api/users", c.UserController.FindAll)
	c.App.Post("/api/users", c.UserController.Register)
	c.App.Put("/api/users/:id", c.UserController.Update)
	c.App.Delete("/api/users/:id", c.UserController.Delete)

	// sales
	c.App.Get("/api/sales", c.SalesController.FindAll)
	c.App.Put("/api/sales/:id", c.SalesController.Update)
	c.App.Delete("/api/sales/:id", c.SalesController.Delete)

	// route
	c.App.Get("/api/routes", c.RouteController.FindAll)
	c.App.Post("/api/routes", c.RouteController.Create)
	c.App.Put("/api/routes/:id", c.RouteController.Update)
	c.App.Delete("/api/routes/:id", c.RouteController.Delete)

	// employee
	c.App.Get("/api/employees", c.EmployeeController.FindAll)
	c.App.Get("/api/employees/:id", c.EmployeeController.FindById)
	c.App.Post("/api/employees", c.EmployeeController.Create)
	c.App.Put("/api/employees/:id", c.EmployeeController.Update)
	c.App.Delete("/api/employees/:id", c.EmployeeController.Delete)

	// Attendace
	c.App.Get("/api/attendance", c.EmployeeController.FindAllWithAttendances)
	c.App.Post("/api/attendance/batch", c.EmployeeAttendanceController.Upsert)

}
