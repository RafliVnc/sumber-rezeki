package http

import (
	"api/internal/model"
	"api/internal/usecase"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type FactoryController struct {
	Log            *logrus.Logger
	FactoryUseCase usecase.FactoryUseCase
}

func NewFactoryController(useCase usecase.FactoryUseCase, logger *logrus.Logger) *FactoryController {
	return &FactoryController{
		FactoryUseCase: useCase,
		Log:            logger,
	}
}

func (c *FactoryController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllFactoryRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Search:  ctx.Query("search"),
	}

	response, total, err := c.FactoryUseCase.FindAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting factories")
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

	return ctx.JSON(model.WebResponse[[]model.FactoryResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *FactoryController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateFactoryRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.FactoryUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create factory : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.FactoryResponse]{Data: response})
}

func (c *FactoryController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateFactoryRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request.ID = int64(id)

	response, err := c.FactoryUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating factory")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.FactoryResponse]{Data: response})
}

func (c *FactoryController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.DeleteFactoryRequest{
		ID: int64(id),
	}

	if err := c.FactoryUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting Factory")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
