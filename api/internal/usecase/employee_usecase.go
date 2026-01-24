package usecase

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"api/internal/model"
	"api/internal/model/converter"
	"api/internal/repository"
	"api/internal/utils"
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeUseCase interface {
	FindAll(ctx context.Context, request *model.FindAllEmployeeRequest) ([]model.EmployeeResponse, int64, error)
	Create(ctx context.Context, request *model.CreateEmployeeRequest) (*model.EmployeeResponse, error)
	Update(ctx context.Context, request *model.UpdateEmployeeRequest) (*model.EmployeeResponse, error)
	Delete(ctx context.Context, request *model.DeleteEmployeeRequest) error
	FindById(ctx context.Context, request *model.FindByIdEmployeeRequest) (*model.EmployeeResponse, error)
	validateAndGetRoutes(tx *gorm.DB, routeIDs []int) ([]entity.Route, error)
	FindAllWithAttendances(ctx context.Context, request *model.FindAllEmployeeWithAttendanceRequest) ([]model.EmployeeResponse, error)
	// TODO FindEmployeeAttendace
}

type EmployeeUseCaseImpl struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validate           *validator.Validate
	EmployeeRepository repository.EmployeeRepository
	RouteRepository    repository.RouteRepository
	SalesRepository    repository.SalesRepository
}

func NewEmployeeUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	employeeRepository repository.EmployeeRepository, routeRepository repository.RouteRepository,
	salesRepository repository.SalesRepository) EmployeeUseCase {
	return &EmployeeUseCaseImpl{
		DB:                 db,
		Log:                logger,
		Validate:           validate,
		EmployeeRepository: employeeRepository,
		RouteRepository:    routeRepository,
		SalesRepository:    salesRepository,
	}
}

func (u *EmployeeUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllEmployeeRequest) ([]model.EmployeeResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	employees, total, err := u.EmployeeRepository.FindAll(u.DB.WithContext(ctx), request)
	if err != nil {
		u.Log.WithError(err).Error("error getting employees")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
	responses := make([]model.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *converter.ToEmployeeResponse(&employee)
	}

	return responses, total, nil
}

func (u *EmployeeUseCaseImpl) Create(ctx context.Context, request *model.CreateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Create new employee FIRST
	newEmployee := &entity.Employee{
		Name:     request.Name,
		Salary:   request.Salary,
		Role:     request.Role,
		JoinDate: time.Now(),
	}

	// Set supervisor for driver/helper
	if request.Role == enum.DRIVER || request.Role == enum.HELPER {
		if request.SupervisorId != 0 {
			newEmployee.SupervisorId = &request.SupervisorId
		}
	}

	// Save Employee first to get ID
	if err := u.EmployeeRepository.Create(tx, newEmployee); err != nil {
		u.Log.Warnf("Failed create Employee to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Create Sales if role is SALES (AFTER employee created)
	if request.Role == enum.SALES {
		var routes []entity.Route

		// Validate routes if provided
		if len(request.RouteIDs) > 0 {
			dbRoutes, err := u.RouteRepository.FindByArryId(tx, request.RouteIDs)
			if err != nil {
				u.Log.Warnf("Failed find route to database : %+v", err)
				return nil, fiber.ErrInternalServerError
			}

			foundRoutes := make(map[uint]bool)
			for _, route := range dbRoutes {
				foundRoutes[uint(route.ID)] = true
			}

			// Check if all routes exist
			var missingRoutes []int
			for _, requestedID := range request.RouteIDs {
				if !foundRoutes[uint(requestedID)] {
					missingRoutes = append(missingRoutes, requestedID)
				}
			}

			if len(missingRoutes) > 0 {
				u.Log.Warnf("Routes not found : %v", missingRoutes)
				errorMessage := fmt.Sprintf("Route dengan ID %v tidak ditemukan", missingRoutes)
				return nil, fiber.NewError(fiber.StatusNotFound, errorMessage)
			}

			// Set routes for association
			routes = make([]entity.Route, len(request.RouteIDs))
			for i, v := range request.RouteIDs {
				routes[i] = entity.Route{ID: v}
			}
		}

		// Create Sales with EmployeeId reference
		sales := &entity.Sales{
			EmployeeId: newEmployee.ID,
			Phone:      request.Phone,
			Routes:     routes,
		}

		if err := u.SalesRepository.Create(tx, sales); err != nil {
			u.Log.Warnf("Failed create sales to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		u.Log.WithFields(logrus.Fields{
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToEmployeeResponse(newEmployee), nil
}

func (u *EmployeeUseCaseImpl) Update(ctx context.Context, request *model.UpdateEmployeeRequest) (*model.EmployeeResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	// Validate request
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Find existing employee
	employee, err := u.EmployeeRepository.FindByIdWithSubordinates(tx, request.ID)
	if err != nil {
		u.Log.Warnf("Failed find employee to database : %+v", err)
		return nil, fiber.ErrNotFound
	}

	if employee == nil {
		u.Log.Warnf("Employee not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Pegawai tidak ditemukan")
	}

	oldRole := employee.Role
	newRole := request.Role

	// Update basic employee data
	employee.Name = request.Name
	employee.Salary = request.Salary
	employee.Role = request.Role

	// Handle supervisor for driver/helper
	if request.Role == enum.DRIVER || request.Role == enum.HELPER {
		if request.SupervisorId != 0 {
			employee.SupervisorId = &request.SupervisorId
		} else {
			employee.SupervisorId = nil
		}
	} else {
		employee.SupervisorId = nil
	}

	// Update employee
	if err := u.EmployeeRepository.Update(tx, employee); err != nil {
		u.Log.Warnf("Failed update Employee to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	// Handle role change scenarios
	switch {
	case (oldRole == enum.DRIVER || oldRole == enum.HELPER) && (newRole != enum.DRIVER && newRole != enum.HELPER):
		// Driver/Helper → Non-Driver/Helper: Clear supervisor
		employee.SupervisorId = nil

	case oldRole == enum.SALES && newRole != enum.SALES:
		// SALES → Non-SALES: Delete sales record
		if len(employee.Subordinates) > 0 {
			u.Log.Warnf("Employee has subordinates : %d", request.ID)
			return nil, fiber.NewError(fiber.StatusBadRequest, "Sales masih memiliki bawahan")
		}

		if err := u.SalesRepository.DeleteByEmployeeId(tx, employee.ID); err != nil {
			u.Log.Warnf("Failed delete sales record : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

	case oldRole != enum.SALES && newRole == enum.SALES:
		// Non-SALES → SALES: Create new sales record
		var routes []entity.Route
		if request.RouteIDs != nil && len(*request.RouteIDs) > 0 {
			routes, err = u.validateAndGetRoutes(tx, *request.RouteIDs)
			if err != nil {
				return nil, err
			}
		}

		sales := &entity.Sales{
			EmployeeId: employee.ID,
			Phone:      request.Phone,
			Routes:     routes,
		}

		if err := u.SalesRepository.Create(tx, sales); err != nil {
			u.Log.Warnf("Failed create sales to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

	case oldRole == enum.SALES && newRole == enum.SALES:
		// SALES → SALES: Update existing sales record
		sales, err := u.SalesRepository.FindByEmployeeId(tx, employee.ID)
		if err != nil {
			u.Log.Warnf("Failed find sales record : %+v", err)
			return nil, fiber.ErrInternalServerError
		}

		if sales == nil {
			u.Log.Warnf("Sales not found : %d", employee.ID)
			return nil, fiber.NewError(fiber.StatusNotFound, "Sales tidak ditemukan")
		}

		// Update phone
		sales.Phone = request.Phone

		// Update routes if provided
		if request.RouteIDs != nil {
			routes, err := u.validateAndGetRoutes(tx, *request.RouteIDs)
			if err != nil {
				return nil, err
			}

			// Replace routes using GORM Association
			if err := u.SalesRepository.ReplaceRoutes(tx, sales, routes); err != nil {
				u.Log.Warnf("Failed replace routes to database : %+v", err)
				return nil, fiber.ErrInternalServerError
			}
		}

		// Update sales record
		if err := u.SalesRepository.Update(tx, sales); err != nil {
			u.Log.Warnf("Failed update sales to database : %+v", err)
			return nil, fiber.ErrInternalServerError
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		u.Log.WithFields(logrus.Fields{
			"id":   request.ID,
			"name": request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToEmployeeResponse(employee), nil
}

func (u *EmployeeUseCaseImpl) validateAndGetRoutes(tx *gorm.DB, routeIDs []int) ([]entity.Route, error) {
	dbRoutes, err := u.RouteRepository.FindByArryId(tx, routeIDs)
	if err != nil {
		u.Log.Warnf("Failed find route to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	foundRoutes := make(map[uint]bool)
	for _, route := range dbRoutes {
		foundRoutes[uint(route.ID)] = true
	}

	// Check if all routes exist
	var missingRoutes []int
	for _, requestedID := range routeIDs {
		if !foundRoutes[uint(requestedID)] {
			missingRoutes = append(missingRoutes, requestedID)
		}
	}

	if len(missingRoutes) > 0 {
		u.Log.Warnf("Routes not found : %v", missingRoutes)
		errorMessage := fmt.Sprintf("Route dengan ID %v tidak ditemukan", missingRoutes)
		return nil, fiber.NewError(fiber.StatusNotFound, errorMessage)
	}

	// Create routes for association
	routes := make([]entity.Route, len(routeIDs))
	for i, v := range routeIDs {
		routes[i] = entity.Route{ID: v}
	}

	return routes, nil
}

func (u *EmployeeUseCaseImpl) Delete(ctx context.Context, request *model.DeleteEmployeeRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	// Check if employee exists
	dbEmployee, err := u.EmployeeRepository.FindByIdWithSubordinates(tx, request.ID)
	if err != nil {
		u.Log.Warnf("Failed find employee to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	if dbEmployee == nil {
		u.Log.Warnf("Employee not found : %d", request.ID)
		return fiber.NewError(fiber.StatusNotFound, "Karyawan tidak ditemukan")
	}

	// check if have Subordinates
	if len(dbEmployee.Subordinates) > 0 {
		u.Log.Warnf("Employee has subordinates : %d", request.ID)
		return fiber.NewError(fiber.StatusBadRequest, "Sales masih memiliki bawahan")
	}

	// check if sales
	if dbEmployee.Role == enum.SALES {
		if err := u.SalesRepository.DeleteByEmployeeId(tx, dbEmployee.ID); err != nil {
			u.Log.Warnf("Failed delete sales record : %+v", err)
			return fiber.ErrInternalServerError
		}
	}

	if err := u.EmployeeRepository.Delete(tx, request.ID); err != nil {
		u.Log.WithError(err).Error("error deleting employee")
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

func (u *EmployeeUseCaseImpl) FindById(ctx context.Context, request *model.FindByIdEmployeeRequest) (*model.EmployeeResponse, error) {

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	employee, err := u.EmployeeRepository.FindById(u.DB.WithContext(ctx), request.ID)
	if err != nil {
		u.Log.WithError(err).Error("error getting employee")
		return nil, fiber.ErrInternalServerError
	}

	if employee == nil {
		u.Log.Warnf("Employee not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Karyawan tidak ditemukan")
	}

	response := converter.ToEmployeeResponse(employee)

	if employee.Sales != nil {
		response.Sales = &model.SalesResponse{
			ID:        employee.Sales.ID,
			Phone:     employee.Sales.Phone,
			CreatedAt: employee.Sales.CreatedAt,
		}

		if len(employee.Sales.Routes) > 0 {
			response.Sales.Routes = make([]model.RouteResponse, len(employee.Sales.Routes))
			for i, route := range employee.Sales.Routes {
				response.Sales.Routes[i] = model.RouteResponse{
					ID:   route.ID,
					Name: route.Name,
				}
			}
		}
	}

	return response, nil
}

func (u *EmployeeUseCaseImpl) FindAllWithAttendances(ctx context.Context, request *model.FindAllEmployeeWithAttendanceRequest) ([]model.EmployeeResponse, error) {

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		u.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	employees, err := u.EmployeeRepository.FindAllWithAttendances(u.DB.WithContext(ctx), request)
	if err != nil {
		u.Log.WithError(err).Error("error getting employees")
		return nil, fiber.ErrInternalServerError
	}

	responses := make([]model.EmployeeResponse, len(employees))
	for i, employee := range employees {
		responses[i] = *converter.ToEmployeeResponse(&employee)

		// Convert attendances
		attendances := make([]model.EmployeeAttendanceResponse, len(employee.EmployeeAttendance))
		for j, attendance := range employee.EmployeeAttendance {
			attendances[j] = *converter.ToEmployeeAttendanceResponse(&attendance)
		}
		responses[i].Attendaces = attendances
	}

	return responses, nil
}
