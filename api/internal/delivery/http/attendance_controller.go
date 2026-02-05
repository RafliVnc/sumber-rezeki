package http

import (
	"api/internal/model"
	"api/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeAttendanceController struct {
	Log                       *logrus.Logger
	EmployeeAttendanceUseCase usecase.EmployeeAttendanceUseCase
}

func NewEmployeeAttendanceController(useCase usecase.EmployeeAttendanceUseCase, logger *logrus.Logger) *EmployeeAttendanceController {
	return &EmployeeAttendanceController{
		EmployeeAttendanceUseCase: useCase,
		Log:                       logger,
	}
}

func (c *EmployeeAttendanceController) Upsert(ctx *fiber.Ctx) error {
	request := new(model.UpsertEmployeeAttendanceRequest)

	// Parse body
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	newError := c.EmployeeAttendanceUseCase.Upsert(ctx.UserContext(), request)
	if newError != nil {
		c.Log.Warnf("Failed to upsert employee attendance : %+v", err)
		return newError
	}

	return ctx.Status(fiber.StatusOK).
		JSON(model.WebResponse[interface{}]{Data: true})
}
