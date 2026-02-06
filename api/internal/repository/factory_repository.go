package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FactoryRepository interface {
	FindAll(db *gorm.DB, request *model.FindAllFactoryRequest) ([]entity.Factory, int64, error)
	Create(db *gorm.DB, entity *entity.Factory) error
	Update(db *gorm.DB, id int64, updates any) error
	Delete(db *gorm.DB, id int64) error
	FindById(db *gorm.DB, id int64) (*entity.Factory, error)
	CountByPhone(db *gorm.DB, phone string) (int64, error)
}

type factoryRepositoryImpl struct {
	Log *logrus.Logger
}

func NewFactoryRepository(log *logrus.Logger) FactoryRepository {
	return &factoryRepositoryImpl{
		Log: log}
}

func (r *factoryRepositoryImpl) Create(db *gorm.DB, entity *entity.Factory) error {
	return db.Create(entity).Error
}

func (r *factoryRepositoryImpl) Update(db *gorm.DB, id int64, updates any) error {
	return db.Model(&entity.Factory{}).Where("id = ?", id).Updates(updates).Error
}

func (r *factoryRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllFactoryRequest) ([]entity.Factory, int64, error) {
	var factories []entity.Factory
	var total int64

	countQuery := db.Model(new(entity.Factory)).Scopes(r.FilterFactory(request))
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("failed to count factories")
		return nil, 0, err
	}

	query := db.Model(new(entity.Factory)).Scopes(r.FilterFactory(request)).Order("name DESC")

	if request.Page > 0 && request.PerPage > 0 {
		offset := (request.Page - 1) * request.PerPage
		query = query.Offset(offset).Limit(request.PerPage)
	}

	if err := query.Find(&factories).Error; err != nil {
		r.Log.WithError(err).Error("failed to find factories")
		return nil, 0, err
	}

	return factories, total, nil
}

func (r *factoryRepositoryImpl) FilterFactory(request *model.FindAllFactoryRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if search := request.Search; search != "" {
			search = "%" + search + "%"
			tx = tx.Where("name ILIKE ? OR phone ILIKE ?", search, search)
		}

		return tx
	}
}

func (r *factoryRepositoryImpl) Delete(db *gorm.DB, id int64) error {
	return db.Delete(&entity.Factory{}, id).Error
}

func (r *factoryRepositoryImpl) FindById(db *gorm.DB, id int64) (*entity.Factory, error) {
	factory := &entity.Factory{}
	if err := db.First(&factory, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return factory, nil
}

func (r *factoryRepositoryImpl) CountByPhone(db *gorm.DB, phone string) (int64, error) {
	var count int64

	// Cek with soft delete
	err := db.Model(new(entity.Factory)).Where("phone = ?", phone).Count(&count).Error
	return count, err
}
