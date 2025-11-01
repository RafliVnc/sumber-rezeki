package http

import (
	"api/internal/model"
	"api/internal/usecase"
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type RouteController struct {
	Log          *logrus.Logger
	RouteUseCase usecase.RouteUseCase
}

func NewRouteController(useCase usecase.RouteUseCase, logger *logrus.Logger) *RouteController {
	return &RouteController{
		RouteUseCase: useCase,
		Log:          logger,
	}
}

func (c *RouteController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllRouteRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Search:  ctx.Query("search"),
	}

	response, total, err := c.RouteUseCase.FindAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting routes")
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

	return ctx.JSON(model.WebResponse[[]model.RouteResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *RouteController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateRouteRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.RouteUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create route : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.RouteResponse]{Data: response})
}

func (c *RouteController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateRouteRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request.ID = id

	response, err := c.RouteUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating route")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.RouteResponse]{Data: response})
}

func (c *RouteController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.DeleteRouteRequest{
		ID: id,
	}

	if err := c.RouteUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting Route")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
