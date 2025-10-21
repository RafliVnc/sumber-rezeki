package config

import (
	"api/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ErrorHandler: NewErrorHandler(),
		Prefork:      config.GetBool("web.prefork"),
	})

	return app
}

func NewErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		// Check if it's a custom ErrorResponse
		if e, ok := err.(*model.ErrorResponse); ok {
			return ctx.Status(e.Code).JSON(e)
		}

		// Check if it's a fiber.Error
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		// Default error response
		errorMessage := &model.ErrorResponse{
			Code:    code,
			Message: err.Error(),
		}

		return ctx.Status(code).JSON(errorMessage)
	}
}
