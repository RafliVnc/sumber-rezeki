package config

import (
	"api/internal/delivery/http"
	"api/internal/delivery/http/middleware"
	"api/internal/delivery/http/route"
	"api/internal/repository"
	"api/internal/usecase"
	"api/internal/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Log      *logrus.Logger
	App      *fiber.App
	DB       *gorm.DB
	Config   *viper.Viper
	Validate *validator.Validate
	Redis    *redis.Client
}

func Bootstrap(config *BootstrapConfig) {
	utils.InitValidator()
	tokenUtil := utils.NewTokenUtil(config.Config.GetString("secret_key"), config.Redis)

	// Repository
	userRepository := repository.NewUserRepository(config.Log)
	routeRepository := repository.NewRouteRepository(config.Log)
	salesRepository := repository.NewSalesRepository(config.Log)
	employeeRepository := repository.NewEmployeeRepository(config.Log)

	// UseCase
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenUtil)
	routeUseCase := usecase.NewRouteUseCase(config.DB, config.Log, config.Validate, routeRepository, routeRepository)
	salesUseCase := usecase.NewSalesUseCase(config.DB, config.Log, config.Validate, salesRepository, routeRepository)
	employeeUseCase := usecase.NewEmployeeUseCase(config.DB, config.Log, config.Validate, employeeRepository, routeRepository, salesRepository)

	// Controller
	userController := http.NewUserController(userUseCase, config.Log)
	routeController := http.NewRouteController(routeUseCase, config.Log)
	salesController := http.NewSalesController(salesUseCase, config.Log)
	employeeController := http.NewEmployeeController(employeeUseCase, config.Log)

	// hello
	helloController := http.NewHelloController()

	authMiddleware := middleware.NewAuth(tokenUtil, config.Log, config.Redis)

	routeConfig := route.RouteConfig{
		App:                config.App,
		Config:             config.Config,
		UserController:     userController,
		SalesController:    salesController,
		RouteController:    routeController,
		EmployeeController: employeeController,
		HelloController:    helloController,
		AuthMiddleware:     authMiddleware,
	}

	routeConfig.Setup()
}
