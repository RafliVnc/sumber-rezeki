package config

import (
	"api/internal/delivery/http"
	"api/internal/delivery/http/route"
	"api/internal/repository"
	"api/internal/usecase"

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

	//user
	userRepository := repository.NewUserRepository(config.Log)
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	userController := http.NewUserController(userUseCase, config.Log)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
	}

	routeConfig.Setup()
}
