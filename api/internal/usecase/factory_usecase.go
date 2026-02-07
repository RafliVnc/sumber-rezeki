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

type FactoryUseCase interface {
	Create(ctx context.Context, request *model.CreateFactoryRequest) (*model.FactoryResponse, error)
	FindAll(ctx context.Context, request *model.FindAllFactoryRequest) ([]model.FactoryResponse, int64, error)
	Update(ctx context.Context, request *model.UpdateFactoryRequest) (*model.FactoryResponse, error)
	Delete(ctx context.Context, request *model.DeleteFactoryRequest) error
}

type FactoryUseCaseImpl struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	FactoryRepository repository.FactoryRepository
}

func NewFactoryUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	factoryRepository repository.FactoryRepository,
) FactoryUseCase {
	return &FactoryUseCaseImpl{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		FactoryRepository: factoryRepository,
	}
}

// Helper fuction
func (s *FactoryUseCaseImpl) validateFactoryExists(tx *gorm.DB, id int64) (*entity.Factory, error) {
	factory, err := s.FactoryRepository.FindById(tx, id)
	if err != nil {
		s.Log.Warnf("Failed find factory to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if factory == nil {
		s.Log.Warnf("Factory not found : %d", id)
		return nil, fiber.NewError(fiber.StatusNotFound, "Kilang tidak ditemukan")
	}

	return factory, nil
}

func (s *FactoryUseCaseImpl) validatePhoneUniqueness(tx *gorm.DB, phone string, excludePhone string) error {
	phoneCount, err := s.FactoryRepository.CountByPhone(tx, phone)
	if err != nil {
		s.Log.Warnf("Failed check phone to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if excludePhone != "" && excludePhone == phone {
		return nil
	}

	if phoneCount > 0 {
		s.Log.Warnf("Phone already exists : %s", phone)
		errorMessage := fmt.Sprintf("Nomor HP %s sudah digunakan", phone)
		return fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	return nil
}

// Usecase
func (s *FactoryUseCaseImpl) Create(ctx context.Context, request *model.CreateFactoryRequest) (*model.FactoryResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil || details != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check phone uniqueness
	if err := s.validatePhoneUniqueness(tx, request.Phone, ""); err != nil {
		return nil, err
	}

	//set factory
	factory := &entity.Factory{
		Name:        request.Name,
		Phone:       request.Phone,
		DueDate:     request.DueDate,
		Description: request.Description,
	}

	if err := s.FactoryRepository.Create(tx, factory); err != nil {
		s.Log.Warnf("Failed create factory to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToFactoryResponse(factory), nil

}

func (s *FactoryUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllFactoryRequest) ([]model.FactoryResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get factories
	factories, total, err := s.FactoryRepository.FindAll(s.DB.WithContext(ctx), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting factories")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
	responses := make([]model.FactoryResponse, len(factories))
	for i, factory := range factories {
		responses[i] = *converter.ToFactoryResponse(&factory)
	}

	return responses, total, nil
}

func (s *FactoryUseCaseImpl) Update(ctx context.Context, request *model.UpdateFactoryRequest) (*model.FactoryResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Check if factory exists
	factory, err := s.validateFactoryExists(tx, request.ID)
	if err != nil {
		return nil, err
	}

	// Check phone uniqueness
	if err := s.validatePhoneUniqueness(tx, request.Phone, factory.Phone); err != nil {
		return nil, err
	}

	//set factory
	updateFactory := &entity.Factory{
		ID:          factory.ID,
		Name:        request.Name,
		Phone:       request.Phone,
		DueDate:     request.DueDate,
		Description: request.Description,
	}

	if err := s.FactoryRepository.Update(tx, factory.ID, updateFactory); err != nil {
		s.Log.Warnf("Failed update factory to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToFactoryResponse(updateFactory), nil
}

func (s *FactoryUseCaseImpl) Delete(ctx context.Context, request *model.DeleteFactoryRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check if factory exists
	if _, err := s.validateFactoryExists(tx, request.ID); err != nil {
		return err
	}

	// Delete factory
	if err := s.FactoryRepository.Delete(tx, request.ID); err != nil {
		s.Log.WithError(err).Error("error deleting factory")
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
