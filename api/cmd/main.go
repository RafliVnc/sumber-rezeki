package main

import (
	"api/internal/config"
	"fmt"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	app := config.NewFiber(viperConfig)
	db := config.NewDatabase(viperConfig, log, false)
	validator := config.NewValidator(viperConfig)
	redis := config.NewRedis(viperConfig)

	//app config
	config.Bootstrap(&config.BootstrapConfig{
		Log:      log,
		App:      app,
		DB:       db,
		Config:   viperConfig,
		Validate: validator,
		Redis:    redis,
	})

	webPort := viperConfig.GetInt("web.port")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
