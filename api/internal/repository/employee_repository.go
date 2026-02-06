package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmployeeRepository interface {
	FindAll(db *gorm.DB, request *model.FindAllEmployeeRequest) ([]entity.Employee, int64, error)
	Create(db *gorm.DB, employee *entity.Employee) error
	Update(db *gorm.DB, employee *entity.Employee) error
	Delete(db *gorm.DB, id int) error
	FindById(db *gorm.DB, id int) (*entity.Employee, error)
	FindByIdWithSubordinates(db *gorm.DB, id int) (*entity.Employee, error)
	FindAllWithAttendances(db *gorm.DB, request *model.FindAllEmployeeWithAttendanceRequest) ([]entity.Employee, error)
}

type employeeRepositoryImpl struct {
	Log *logrus.Logger
}

func NewEmployeeRepository(log *logrus.Logger) EmployeeRepository {
	return &employeeRepositoryImpl{
		Log: log,
	}
}

func (r *employeeRepositoryImpl) Create(db *gorm.DB, employee *entity.Employee) error {
	return db.Create(employee).Error
}

func (r *employeeRepositoryImpl) Update(db *gorm.DB, employee *entity.Employee) error {
	return db.Model(&entity.Employee{ID: employee.ID}).Updates(employee).Error
}

func (r *employeeRepositoryImpl) Delete(db *gorm.DB, id int) error {
	return db.Delete(&entity.Employee{}, id).Error
}

func (r *employeeRepositoryImpl) FindById(db *gorm.DB, id int) (*entity.Employee, error) {
	var employee entity.Employee

	err := db.Preload("Sales.Routes").First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllEmployeeRequest) ([]entity.Employee, int64, error) {
	var employees []entity.Employee
	var total int64

	countQuery := db.Model(new(entity.Employee)).Scopes(r.FilterEmployee(request))
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("failed to count employees")
		return nil, 0, err
	}

	query := db.Model(new(entity.Employee)).Scopes(r.FilterEmployee(request)).Order("name ASC")

	if request.Page > 0 && request.PerPage > 0 {
		offset := (request.Page - 1) * request.PerPage
		query = query.Offset(offset).Limit(request.PerPage)
	}

	if err := query.Find(&employees).Error; err != nil {
		r.Log.WithError(err).Error("failed to find employees")
		return nil, 0, err
	}

	return employees, total, nil
}

func (r *employeeRepositoryImpl) FilterEmployee(request *model.FindAllEmployeeRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name ILIKE ? ", name)
		}

		if name := request.Salary; name > 0 {
			tx = tx.Where("salary = ?", name)
		}

		if len(request.Roles) > 0 {
			tx = tx.Where("role IN ?", request.Roles)
		}

		return tx
	}
}

func (r *employeeRepositoryImpl) FindByIdWithSubordinates(db *gorm.DB, id int) (*entity.Employee, error) {
	var employee entity.Employee

	err := db.Preload("Subordinates").First(&employee, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &employee, nil
}

func (r *employeeRepositoryImpl) FindAllWithAttendances(db *gorm.DB, request *model.FindAllEmployeeWithAttendanceRequest) ([]entity.Employee, error) {
	var employees []entity.Employee

	query := db.Preload("EmployeeAttendance", func(db *gorm.DB) *gorm.DB {
		if request.StartDate != "" && request.EndDate != "" {
			return db.Where("date BETWEEN ? AND ?", request.StartDate, request.EndDate)
		}
		return db
	})

	if err := query.Where("join_date <= ?", request.EndDate).
		Find(&employees).Error; err != nil {
		r.Log.WithError(err).Error("failed to find employees")
		return nil, err
	}

	return employees, nil
}
