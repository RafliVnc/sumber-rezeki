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
	FactoryController            *http.FactoryController
	VehicleController            *http.VehicleController
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

	c.App.Get("/api/current", c.UserController.Current)
	// hello
	c.App.Get("/api/hello", c.HelloController.SayHello)

	// user
	users := c.App.Group("/api/users")
	users.Get("/", c.UserController.FindAll)
	users.Post("/", c.UserController.Register)
	users.Put("/:id", c.UserController.Update)
	users.Delete("/:id", c.UserController.Delete)

	// sales
	sales := c.App.Group("/api/sales")
	sales.Get("/", c.SalesController.FindAll)
	sales.Put("/:id", c.SalesController.Update)
	sales.Delete("/:id", c.SalesController.Delete)

	// route
	routes := c.App.Group("/api/routes")
	routes.Get("/", c.RouteController.FindAll)
	routes.Post("/", c.RouteController.Create)
	routes.Put("/:id", c.RouteController.Update)
	routes.Delete("/:id", c.RouteController.Delete)

	// employee
	employees := c.App.Group("/api/employees")
	employees.Get("/", c.EmployeeController.FindAll)
	employees.Get("/:id", c.EmployeeController.FindById)
	employees.Post("/", c.EmployeeController.Create)
	employees.Put("/:id", c.EmployeeController.Update)
	employees.Delete("/:id", c.EmployeeController.Delete)

	// attendance
	attendance := c.App.Group("/api/attendance")
	attendance.Get("/", c.EmployeeController.FindAllWithAttendances)
	attendance.Post("/batch", c.EmployeeAttendanceController.Upsert)

	// factory
	factories := c.App.Group("/api/factories")
	factories.Get("/", c.FactoryController.FindAll)
	factories.Post("/", c.FactoryController.Create)
	factories.Put("/:id", c.FactoryController.Update)
	factories.Delete("/:id", c.FactoryController.Delete)

	// factory
	vehicles := c.App.Group("/api/vehicles")
	vehicles.Get("/", c.VehicleController.FindAll)
	vehicles.Post("/", c.VehicleController.Create)
	vehicles.Put("/:id", c.VehicleController.Update)
	vehicles.Delete("/:id", c.VehicleController.Delete)
}
