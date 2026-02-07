package usecase

import (
	"api/internal/entity"
	"api/internal/model"
	"api/internal/model/converter"
	"api/internal/repository"
	"api/internal/utils"
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VehicleUseCase interface {
	Create(ctx context.Context, request *model.CreateVehicleRequest) (*model.VehicleResponse, error)
	FindAll(ctx context.Context, request *model.FindAllVehicleRequest) ([]model.VehicleResponse, int64, error)
	Update(ctx context.Context, request *model.UpdateVehicleRequest) (*model.VehicleResponse, error)
	Delete(ctx context.Context, request *model.DeleteVehicleRequest) error
}

type VehicleUseCaseImpl struct {
	DB                *gorm.DB
	Log               *logrus.Logger
	Validate          *validator.Validate
	VehicleRepository repository.VehicleRepository
}

func NewVehicleUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	vehicleRepository repository.VehicleRepository,
) VehicleUseCase {
	return &VehicleUseCaseImpl{
		DB:                db,
		Log:               logger,
		Validate:          validate,
		VehicleRepository: vehicleRepository,
	}
}

// Helper fuction
func (s *VehicleUseCaseImpl) validateVehicleExists(tx *gorm.DB, id int64) (*entity.Vehicle, error) {
	vehicle, err := s.VehicleRepository.FindById(tx, id)
	if err != nil {
		s.Log.Warnf("Failed find vehicle to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if vehicle == nil {
		s.Log.Warnf("Vehicle not found : %d", id)
		return nil, fiber.NewError(fiber.StatusNotFound, "Kendaraan tidak ditemukan")
	}

	return vehicle, nil
}

func (s *VehicleUseCaseImpl) validatePlateUniqueness(tx *gorm.DB, plate string, excludePlate string) error {
	plateCount, err := s.VehicleRepository.CountByPlate(tx, plate)
	if err != nil {
		s.Log.Warnf("Failed check plate to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if excludePlate != "" && excludePlate == plate {
		return nil
	}

	if plateCount > 0 {
		s.Log.Warnf("Plate already exists : %s", plate)
		errorMessage := fmt.Sprintf("Plat Kendaraan %s sudah digunakan", plate)
		return fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	return nil
}

// Usecase
func (s *VehicleUseCaseImpl) Create(ctx context.Context, request *model.CreateVehicleRequest) (*model.VehicleResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil || details != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check plate uniqueness
	if err := s.validatePlateUniqueness(tx, request.Plate, ""); err != nil {
		return nil, err
	}

	//set vehicle
	vehicle := &entity.Vehicle{
		Plate: strings.ToUpper(strings.Join(strings.Fields(request.Plate), "")),
		Type:  request.Type,
	}

	if err := s.VehicleRepository.Create(tx, vehicle); err != nil {
		s.Log.Warnf("Failed create vehicle to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Plate,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToVehicleResponse(vehicle), nil

}

func (s *VehicleUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllVehicleRequest) ([]model.VehicleResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//get vehicles
	vehicles, total, err := s.VehicleRepository.FindAll(s.DB.WithContext(ctx), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting vehicles")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
	responses := make([]model.VehicleResponse, len(vehicles))
	for i, vehicle := range vehicles {
		responses[i] = *converter.ToVehicleResponse(&vehicle)
	}

	return responses, total, nil
}

func (s *VehicleUseCaseImpl) Update(ctx context.Context, request *model.UpdateVehicleRequest) (*model.VehicleResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Check if vehicle exists
	vehicle, err := s.validateVehicleExists(tx, request.ID)
	if err != nil {
		return nil, err
	}

	// Check plate uniqueness
	if err := s.validatePlateUniqueness(tx, request.Plate, vehicle.Plate); err != nil {
		return nil, err
	}

	//set vehicle
	updateVehicle := &entity.Vehicle{
		ID:    vehicle.ID,
		Plate: strings.ToUpper(strings.Join(strings.Fields(request.Plate), "")),
		Type:  request.Type,
	}

	if err := s.VehicleRepository.Update(tx, vehicle.ID, updateVehicle); err != nil {
		s.Log.Warnf("Failed update vehicle to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"name": request.Plate,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToVehicleResponse(updateVehicle), nil
}

func (s *VehicleUseCaseImpl) Delete(ctx context.Context, request *model.DeleteVehicleRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check if vehicle exists
	if _, err := s.validateVehicleExists(tx, request.ID); err != nil {
		return err
	}

	// Delete vehicle
	if err := s.VehicleRepository.Delete(tx, request.ID); err != nil {
		s.Log.WithError(err).Error("error deleting vehicle")
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
