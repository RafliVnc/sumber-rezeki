package usecase

import (
	"api/internal/entity"
	"api/internal/model"
	"api/internal/model/converter"
	"api/internal/repository"
	"api/internal/utils"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RouteUseCase interface {
	FindAll(ctx context.Context, request *model.FindAllRouteRequest) ([]model.RouteResponse, int64, error)
	Create(ctx context.Context, request *model.CreateRouteRequest) (*model.RouteResponse, error)
	Update(ctx context.Context, request *model.UpdateRouteRequest) (*model.RouteResponse, error)
	Delete(ctx context.Context, request *model.DeleteRouteRequest) error
}

type RouteUseCaseImpl struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	RouteRepository repository.RouteRepository
}

func NewRouteUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	routesRepository repository.RouteRepository,
	routeRepository repository.RouteRepository,
) RouteUseCase {
	return &RouteUseCaseImpl{
		DB:              db,
		Log:             logger,
		Validate:        validate,
		RouteRepository: routeRepository,
	}
}

func (u *RouteUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllRouteRequest) ([]model.RouteResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get routes
	routes, total, err := u.RouteRepository.FindAll(u.DB.WithContext(ctx), request)
	if err != nil {
		u.Log.WithError(err).Error("error getting routes")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
	responses := make([]model.RouteResponse, len(routes))
	for i, routes := range routes {
		responses[i] = *converter.ToRouteResponse(&routes)
	}

	return responses, total, nil
}

func (u *RouteUseCaseImpl) Create(ctx context.Context, request *model.CreateRouteRequest) (*model.RouteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check duplicate name
	count, err := u.RouteRepository.CountByName(tx, request.Name)
	if err != nil {
		u.Log.Warnf("Failed check name to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		u.Log.Warnf("Name already exists : %s", request.Name)
		errorMessage := fmt.Sprintf("Nama %s sudah digunakan", request.Name)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	// set route
	route := &entity.Route{
		Name:        request.Name,
		Description: request.Description,
	}

	//create route
	if err := u.RouteRepository.Create(tx, route); err != nil {
		u.Log.Warnf("Failed create route to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		u.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToRouteResponse(route), nil
}

func (u *RouteUseCaseImpl) Update(ctx context.Context, request *model.UpdateRouteRequest) (*model.RouteResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// check db exists
	DbRoute, err := u.RouteRepository.FindById(tx, request.ID)
	if err != nil {
		u.Log.Warnf("Failed find route to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if DbRoute == nil {
		u.Log.Warnf("Route not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Rute tidak ditemukan")
	}

	// Check duplicate name
	count, err := u.RouteRepository.CountByName(tx, request.Name)
	if err != nil {
		u.Log.Warnf("Failed check name to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 && DbRoute.Name != request.Name {
		u.Log.Warnf("Name already exists : %s", request.Name)
		errorMessage := fmt.Sprintf("Nama %s sudah digunakan", request.Name)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	// set route
	route := &entity.Route{
		ID:          request.ID,
		Name:        request.Name,
		Description: request.Description,
	}

	//update route
	if err := u.RouteRepository.Update(tx, route); err != nil {
		u.Log.Warnf("Failed update route to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		u.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToRouteResponse(route), nil
}

func (u *RouteUseCaseImpl) Delete(ctx context.Context, request *model.DeleteRouteRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// check db exists
	DbRoute, err := u.RouteRepository.FindById(tx, request.ID)
	if err != nil {
		u.Log.Warnf("Failed find route to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if DbRoute == nil {
		u.Log.Warnf("Route not found : %d", request.ID)
		return fiber.NewError(fiber.StatusNotFound, "Rute tidak ditemukan")
	}

	if len(DbRoute.Sales) > 0 {
		u.Log.Warnf("Route has sales : %d", request.ID)
		return fiber.NewError(fiber.StatusBadRequest, "Rute digunakan oleh sales")
	}

	//delete route
	if err := u.RouteRepository.Delete(tx, request.ID); err != nil {
		u.Log.Warnf("Failed delete route to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		u.Log.WithFields(logrus.Fields{
			"id": request.ID,
		}).Warnf("Failed commit to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
