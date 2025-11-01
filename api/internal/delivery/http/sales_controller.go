package http

import (
	"api/internal/model"
	"api/internal/usecase"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SalesController struct {
	Log          *logrus.Logger
	SalesUseCase usecase.SalesUseCase
}

func NewSalesController(useCase usecase.SalesUseCase, logger *logrus.Logger) *SalesController {
	return &SalesController{
		SalesUseCase: useCase,
		Log:          logger,
	}
}

func (c *SalesController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateSalesRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.SalesUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create sales : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.SalesResponse]{Data: response})
}

func (c *SalesController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllSalesRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Search:  ctx.Query("search"),
	}

	response, total, err := c.SalesUseCase.FindAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting sales")
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

	return ctx.JSON(model.WebResponse[[]model.SalesResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *SalesController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.DeleteSalesRequest{
		ID: id,
	}

	if err := c.SalesUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting Sales")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}

func (c *SalesController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateSalesRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request.ID = id

	response, err := c.SalesUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating sales")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.SalesResponse]{Data: response})
}
