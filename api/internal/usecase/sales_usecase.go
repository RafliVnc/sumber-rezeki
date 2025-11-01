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

type SalesUseCase interface {
	Create(ctx context.Context, request *model.CreateSalesRequest) (*model.SalesResponse, error)
	FindAll(ctx context.Context, request *model.FindAllSalesRequest) ([]model.SalesResponse, int64, error)
	Update(ctx context.Context, request *model.UpdateSalesRequest) (*model.SalesResponse, error)
	Delete(ctx context.Context, request *model.DeleteSalesRequest) error
}

type SalesUseCaseImpl struct {
	DB              *gorm.DB
	Log             *logrus.Logger
	Validate        *validator.Validate
	SalesRepository repository.SalesRepository
	RouteRepository repository.RouteRepository
}

func NewSalesUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate, salesRepository repository.SalesRepository, routeRepository repository.RouteRepository) SalesUseCase {
	return &SalesUseCaseImpl{
		DB:              db,
		Log:             logger,
		Validate:        validate,
		SalesRepository: salesRepository,
		RouteRepository: routeRepository,
	}
}

func (s *SalesUseCaseImpl) Create(ctx context.Context, request *model.CreateSalesRequest) (*model.SalesResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil || details != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//check phone uniqueness
	count, err := s.SalesRepository.CountByPhone(tx, request.Phone)
	if err != nil {
		s.Log.Warnf("Failed check phone to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		s.Log.Warnf("Phone already exists : %s", request.Phone)
		errorMessage := fmt.Sprintf("Nomor HP %s sudah digunakan", request.Phone)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	var Routes []entity.Route
	//  validate routes
	if request.RouteIDs != nil {
		DbRoutes, err := s.RouteRepository.FindByArryId(tx, request.RouteIDs)
		if err != nil {
			s.Log.Warnf("Failed find route to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		foundRoutes := make(map[uint]bool)
		for _, route := range DbRoutes {
			foundRoutes[uint(route.ID)] = true
		}

		// check if all routes exist
		var missingRoutes []int
		for _, requestedID := range request.RouteIDs {
			if !foundRoutes[uint(requestedID)] {
				missingRoutes = append(missingRoutes, requestedID)
			}
		}

		if len(missingRoutes) > 0 {
			s.Log.Warnf("Routes not found : %v", missingRoutes)
			errorMessage := fmt.Sprintf("Route dengan ID %v tidak ditemukan", missingRoutes)
			return nil, fiber.NewError(fiber.StatusNotFound, errorMessage)
		}

		// set routes
		Routes = make([]entity.Route, len(request.RouteIDs))
		for i, v := range request.RouteIDs {
			Routes[i] = entity.Route{ID: v}
		}

	}

	//set sales
	sales := &entity.Sales{
		Name:   request.Name,
		Phone:  request.Phone,
		Routes: Routes,
	}

	if err := s.SalesRepository.Create(tx, sales); err != nil {
		s.Log.Warnf("Failed create sales to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToSalesResponse(sales), nil
}

func (s *SalesUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllSalesRequest) ([]model.SalesResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get sales
	sales, total, err := s.SalesRepository.FindAll(s.DB.WithContext(ctx), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting sales")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
	responses := make([]model.SalesResponse, len(sales))
	for i, sales := range sales {
		responses[i] = *converter.ToSalesResponse(&sales)
	}

	return responses, total, nil
}

func (s *SalesUseCaseImpl) Update(ctx context.Context, request *model.UpdateSalesRequest) (*model.SalesResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get sales
	DbSales, err := s.SalesRepository.FindById(tx, request.ID)
	if err != nil {
		s.Log.WithError(err).Error("error getting sales")
		return nil, fiber.ErrInternalServerError
	}

	if DbSales == nil {
		s.Log.Warnf("Sales not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Sales tidak ditemukan")
	}

	//check phone uniqueness
	count, err := s.SalesRepository.CountByPhone(tx, request.Phone)
	if err != nil {
		s.Log.Warnf("Failed check phone to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 && DbSales.Phone != request.Phone {
		s.Log.Warnf("Phone already exists : %s", request.Phone)
		errorMessage := fmt.Sprintf("Nomor HP %s sudah digunakan", request.Phone)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//  validate routes
	if request.RouteIDs != nil {
		DbRoutes, err := s.RouteRepository.FindByArryId(tx, request.RouteIDs)
		if err != nil {
			s.Log.Warnf("Failed find route to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		foundRoutes := make(map[uint]bool)
		for _, route := range DbRoutes {
			foundRoutes[uint(route.ID)] = true
		}

		// check if all routes exist
		var missingRoutes []int
		for _, requestedID := range request.RouteIDs {
			if !foundRoutes[uint(requestedID)] {
				missingRoutes = append(missingRoutes, requestedID)
			}
		}

		if len(missingRoutes) > 0 {
			s.Log.Warnf("Routes not found : %v", missingRoutes)
			errorMessage := fmt.Sprintf("Route dengan ID %v tidak ditemukan", missingRoutes)
			return nil, fiber.NewError(fiber.StatusNotFound, errorMessage)
		}

		// create new routes slice
		newRoutes := make([]entity.Route, len(request.RouteIDs))
		for i, routeID := range request.RouteIDs {
			newRoutes[i] = entity.Route{ID: routeID}
		}

		// Replace routes using GORM Association
		if err := s.SalesRepository.ReplaceRoutes(tx, DbSales, newRoutes); err != nil {
			s.Log.Warnf("Failed replace routes to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	//set sales
	sales := &entity.Sales{
		ID:    request.ID,
		Name:  request.Name,
		Phone: request.Phone,
	}

	if err := s.SalesRepository.Update(tx, sales); err != nil {
		s.Log.Warnf("Failed update sales to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToSalesResponse(sales), nil
}

func (s *SalesUseCaseImpl) Delete(ctx context.Context, request *model.DeleteSalesRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get sales
	DbSales, err := s.SalesRepository.FindById(tx, request.ID)
	if err != nil {
		s.Log.WithError(err).Error("error getting sales")
		return fiber.ErrInternalServerError
	}

	if DbSales == nil {
		s.Log.Warnf("Sales not found : %d", request.ID)
		return fiber.NewError(fiber.StatusNotFound, "Sales tidak ditemukan")
	}

	if err := s.SalesRepository.Delete(tx, request.ID); err != nil {
		s.Log.WithError(err).Error("error deleting Sales")
		return fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"id": request.ID,
		}).Warnf("Failed commit to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
