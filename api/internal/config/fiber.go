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
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}

		var errorMessage *model.ErrorResponse
		errorMessage = &model.ErrorResponse{
			Code:    code,
			Message: err.Error(),
		}

		if e, ok := err.(*model.ErrorResponse); ok {
			errorMessage = &model.ErrorResponse{
				Code:    e.Code,
				Message: e.Message,
				Details: e.Details,
			}
		}

		return ctx.Status(code).JSON(errorMessage)
	}
}
