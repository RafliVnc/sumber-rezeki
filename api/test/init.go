package test

import (
	"api/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var app *fiber.App
var db *gorm.DB
var viperConfig *viper.Viper
var log *logrus.Logger
var validate *validator.Validate
var redisClient *redis.Client

func init() {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	db = config.NewDatabase(viperConfig, log, true)
	redisClient = config.NewRedis(viperConfig)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Config:   viperConfig,
		Validate: validate,
		Redis:    redisClient,
	})
}
