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

type VehicleController struct {
	Log            *logrus.Logger
	VehicleUseCase usecase.VehicleUseCase
}

func NewVehicleController(useCase usecase.VehicleUseCase, logger *logrus.Logger) *VehicleController {
	return &VehicleController{
		VehicleUseCase: useCase,
		Log:            logger,
	}
}

func (c *VehicleController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllVehicleRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Search:  strings.ToUpper(strings.Join(strings.Fields(ctx.Query("search")), "")),
	}

	typesRaw := ctx.Context().QueryArgs().PeekMulti("types[]")
	for _, r := range typesRaw {
		t := strings.TrimSpace(string(r))
		if t != "" {
			request.Types = append(request.Types, enum.VehicleType(t))
		}
	}

	response, total, err := c.VehicleUseCase.FindAll(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting vehicles")
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

	return ctx.JSON(model.WebResponse[[]model.VehicleResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *VehicleController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateVehicleRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	request.Plate = strings.ToUpper(strings.Join(strings.Fields(request.Plate), ""))

	response, err := c.VehicleUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create vehicle : %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(model.WebResponse[*model.VehicleResponse]{Data: response})
}

func (c *VehicleController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateVehicleRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request.ID = int64(id)
	request.Plate = strings.ToUpper(strings.Join(strings.Fields(request.Plate), ""))

	response, err := c.VehicleUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating vehicle")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.VehicleResponse]{Data: response})
}

func (c *VehicleController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id parameter")
	}

	request := &model.DeleteVehicleRequest{
		ID: int64(id),
	}

	if err := c.VehicleUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting Vehicle")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
