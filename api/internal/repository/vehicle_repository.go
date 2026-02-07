package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type VehicleRepository interface {
	FindAll(db *gorm.DB, request *model.FindAllVehicleRequest) ([]entity.Vehicle, int64, error)
	Create(db *gorm.DB, entity *entity.Vehicle) error
	Update(db *gorm.DB, id int64, updates any) error
	Delete(db *gorm.DB, id int64) error
	FindById(db *gorm.DB, id int64) (*entity.Vehicle, error)
	CountByPlate(db *gorm.DB, plate string) (int64, error)
}

type vehicleRepositoryImpl struct {
	Log *logrus.Logger
}

func NewVehicleRepository(log *logrus.Logger) VehicleRepository {
	return &vehicleRepositoryImpl{
		Log: log}
}

func (r *vehicleRepositoryImpl) Create(db *gorm.DB, entity *entity.Vehicle) error {
	return db.Create(entity).Error
}

func (r *vehicleRepositoryImpl) Update(db *gorm.DB, id int64, updates any) error {
	return db.Model(&entity.Vehicle{}).Where("id = ?", id).Updates(updates).Error
}

func (r *vehicleRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllVehicleRequest) ([]entity.Vehicle, int64, error) {
	var vehicles []entity.Vehicle
	var total int64

	countQuery := db.Model(new(entity.Vehicle)).Scopes(r.FilterVehicle(request))
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("failed to count vehicles")
		return nil, 0, err
	}

	query := db.Model(new(entity.Vehicle)).Scopes(r.FilterVehicle(request)).Order("plate DESC")

	if request.Page > 0 && request.PerPage > 0 {
		offset := (request.Page - 1) * request.PerPage
		query = query.Offset(offset).Limit(request.PerPage)
	}

	if err := query.Find(&vehicles).Error; err != nil {
		r.Log.WithError(err).Error("failed to find vehicles")
		return nil, 0, err
	}

	return vehicles, total, nil
}

func (r *vehicleRepositoryImpl) FilterVehicle(request *model.FindAllVehicleRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if search := request.Search; search != "" {
			search = "%" + search + "%"
			tx = tx.Where("plate ILIKE ? ", search)
		}

		if len(request.Types) > 0 {
			tx = tx.Where("type IN ?", request.Types)
		}

		return tx
	}
}

func (r *vehicleRepositoryImpl) Delete(db *gorm.DB, id int64) error {
	return db.Delete(&entity.Vehicle{}, id).Error
}

func (r *vehicleRepositoryImpl) FindById(db *gorm.DB, id int64) (*entity.Vehicle, error) {
	vehicle := &entity.Vehicle{}
	if err := db.First(&vehicle, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return vehicle, nil
}

func (r *vehicleRepositoryImpl) CountByPlate(db *gorm.DB, plate string) (int64, error) {
	var count int64

	err := db.Model(&entity.Vehicle{}).Where("plate = ?", plate).Count(&count).Error
	return count, err
}
