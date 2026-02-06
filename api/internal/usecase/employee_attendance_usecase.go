package usecase

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"api/internal/model"
	"api/internal/repository"
	"api/internal/utils"
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeAttendanceUseCase interface {
	Upsert(ctx context.Context, request *model.UpsertEmployeeAttendanceRequest) error
}

type EmployeeAttendanceUseCaseImpl struct {
	DB                           *gorm.DB
	Log                          *logrus.Logger
	Validate                     *validator.Validate
	EmployeeAttendanceRepository repository.EmployeeAttendanceRepository
	PeriodUsecase                PeriodUseCase
}

func NewEmployeeAttendanceUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	employeeRepository repository.EmployeeAttendanceRepository,
	periodUsecase PeriodUseCase,
) EmployeeAttendanceUseCase {
	return &EmployeeAttendanceUseCaseImpl{
		DB:                           db,
		Log:                          logger,
		Validate:                     validate,
		EmployeeAttendanceRepository: employeeRepository,
		PeriodUsecase:                periodUsecase,
	}
}

func (u *EmployeeAttendanceUseCaseImpl) Upsert(ctx context.Context, request *model.UpsertEmployeeAttendanceRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Declare batch data Upsert and Delete
	var dataUpsert []*entity.EmployeeAttendance
	var dataDelete []time.Time

	// Create data partition
	for _, req := range request.Attendances {
		if req.Action == "update" {
			resultData, err := u.CreateUpsertData(ctx, &req)
			if err != nil {
				u.Log.Warnf("Failed to create array data: %+v", err)
				return fiber.ErrInternalServerError
			}
			dataUpsert = append(dataUpsert, resultData...)
		} else {
			newDate, err := time.Parse("2006-01-02", req.Date)
			if err != nil {
				u.Log.Warnf("Failed to parse date: %+v", err)
				return fiber.ErrInternalServerError
			}

			dataDelete = append(dataDelete, newDate)
		}
	}

	// Batch upsert
	if len(dataUpsert) > 0 {
		if err := u.EmployeeAttendanceRepository.BatchUpsert(tx, dataUpsert); err != nil {
			u.Log.Warnf("Failed create attendance to database : %+v", err)
			return fiber.ErrInternalServerError
		}
	}

	// Batch delete (soft delete)
	if len(dataDelete) > 0 {
		if err := u.EmployeeAttendanceRepository.BatchDeleteByDate(tx, dataDelete); err != nil {
			u.Log.Warnf("Failed to delete Employee Attendance from database: %+v", err)
			return fiber.ErrInternalServerError
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		u.Log.Warnf("Failed to commit transaction: %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func (u *EmployeeAttendanceUseCaseImpl) CreateUpsertData(ctx context.Context, request *model.AttendanceAction) ([]*entity.EmployeeAttendance, error) {
	var attendances []*entity.EmployeeAttendance

	// Parse date string to time.Time
	date, err := time.Parse("2006-01-02", request.Date)
	if err != nil {
		u.Log.Warnf("Failed to parse date: %+v", err)
		return nil, err
	}

	periodId, err := u.PeriodUsecase.GetOrCreatePeriodIdByDate(ctx, request.Date)
	if err != nil {
		u.Log.Warnf("Failed to generate period id: %+v", err)
		return nil, err
	}

	// validate periodClosure

	for _, emp := range request.Employees {
		attendance := &entity.EmployeeAttendance{
			Date:       date,
			Status:     enum.AttendanceStatus(emp.Status),
			EmployeeId: emp.ID,
			PeriodId:   periodId,
		}
		attendances = append(attendances, attendance)
	}

	return attendances, nil
}
