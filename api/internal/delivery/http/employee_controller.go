package http

import (
	"api/internal/entity/enum"
	"api/internal/model"
	"api/internal/usecase"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type EmployeeController struct {
	Log             *logrus.Logger
	EmployeeUseCase usecase.EmployeeUseCase
}

func NewEmployeeController(useCase usecase.EmployeeUseCase, logger *logrus.Logger) *EmployeeController {
	return &EmployeeController{
		EmployeeUseCase: useCase,
		Log:             logger,
	}
}

func (c *EmployeeController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllEmployeeRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Name:    ctx.Query("search"),
		// TODO: Salary
	}

	rolesRaw := ctx.Context().QueryArgs().PeekMulti("roles[]")
	for _, r := range rolesRaw {
		role := strings.TrimSpace(string(r))
		if role != "" {
			request.Roles = append(request.Roles, enum.EmployeeRole(role))
		}
	}

	response, total, err := c.EmployeeUseCase.FindAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting user")
		return err
	}

	var paging *model.PageMetadata
	if request.Page > 0 && request.PerPage > 0 {
		paging = &model.PageMetadata{
			Page:      request.Page,
			PerPage:   request.PerPage,
			TotalItem: total,
			TotalPage: int64(math.Ceil(float64(total) / float64(request.PerPage))),
		}
	}

	return ctx.JSON(model.WebResponse[[]model.EmployeeResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *EmployeeController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateEmployeeRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.EmployeeUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create employee : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.EmployeeResponse]{Data: response})
}

func (c *EmployeeController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateEmployeeRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request.ID = id

	response, err := c.EmployeeUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{Data: response})
}

func (c *EmployeeController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.DeleteEmployeeRequest{
		ID: id,
	}

	if err := c.EmployeeUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting employee")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *EmployeeController) FindById(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.FindByIdEmployeeRequest{
		ID: id,
	}

	response, err := c.EmployeeUseCase.FindById(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting employee")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.EmployeeResponse]{Data: response})
}
