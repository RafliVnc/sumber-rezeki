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
}

func Bootstrap(config *BootstrapConfig) {
	utils.InitValidator()
	tokenUtil := utils.NewTokenUtil("testSecret")

	//user
	userRepository := repository.NewUserRepository(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, tokenUtil)
	userController := http.NewUserController(userUseCase, config.Log)

	authMiddleware := middleware.NewAuth(userUseCase, tokenUtil)

	routeConfig := route.RouteConfig{
		App:            config.App,
		Config:         config.Config,
		UserController: userController,
		AuthMiddleware: authMiddleware,
	}

	routeConfig.Setup()
}
