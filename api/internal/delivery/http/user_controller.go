package http

import (
	"api/internal/model"
	"api/internal/usecase"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log         *logrus.Logger
	UserUseCase usecase.UserUseCase
}

func NewUserController(useCase usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		UserUseCase: useCase,
		Log:         logger,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := c.UserUseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) FindAll(ctx *fiber.Ctx) error {

	request := &model.FindAllUserRequest{
		Page:    ctx.QueryInt("page"),
		PerPage: ctx.QueryInt("perPage"),
		Name:    ctx.Query("name"),
		Email:   ctx.Query("email"),
	}

	response, total, err := c.UserUseCase.FindAll(ctx.UserContext(), request)
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

	return ctx.JSON(model.WebResponse[[]model.UserResponse]{
		Data:   response,
		Paging: paging,
	})
}

func (c *UserController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateUserRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	request.ID, _ = ctx.ParamsInt("id")

	response, err := c.UserUseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating user")
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}

func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		c.Log.WithError(err).Error("error parsing user ID from URL")
		return fiber.NewError(fiber.StatusBadRequest, "invalid user ID")
	}

	request := &model.DeleteUserRequest{
		ID: id,
	}

	if err := c.UserUseCase.Delete(ctx.UserContext(), request); err != nil {
		c.Log.WithError(err).Error("error deleting user")
		return err
	}

	return ctx.JSON(model.WebResponse[bool]{Data: true})
}
