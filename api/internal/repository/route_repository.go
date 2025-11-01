package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RouteRepository interface {
	FindByArryId(db *gorm.DB, ids []int) ([]entity.Route, error)
	FindAll(db *gorm.DB, request *model.FindAllRouteRequest) ([]entity.Route, int64, error)
	Create(db *gorm.DB, route *entity.Route) error
	Update(db *gorm.DB, route *entity.Route) error
	Delete(db *gorm.DB, id int) error
	CountByName(db *gorm.DB, name string) (int64, error)
	FindById(db *gorm.DB, id int) (*entity.Route, error)
}

type routeRepositoryImpl struct {
	Log *logrus.Logger
}

func NewRouteRepository(log *logrus.Logger) RouteRepository {
	return &routeRepositoryImpl{
		Log: log,
	}
}

func (r *routeRepositoryImpl) FindByArryId(db *gorm.DB, ids []int) ([]entity.Route, error) {
	var routes []entity.Route
	if err := db.Where("id in (?)", ids).Find(&routes).Error; err != nil {
		r.Log.WithError(err).Error("error getting routes")
		return nil, err
	}
	return routes, nil
}

func (r *routeRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllRouteRequest) ([]entity.Route, int64, error) {
	var routesList []entity.Route
	var total int64

	// Count total with filters
	countQuery := db.Model(&entity.Route{}).Scopes(r.FilterRoute(request))
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("failed to count routes")
		return nil, 0, err
	}

	// Main query with filters, preload, and pagination
	query := db.Model(&entity.Route{}).
		Scopes(r.FilterRoute(request)).
		Order("name DESC")

	// Pagination
	if request.Page > 0 && request.PerPage > 0 {
		offset := (request.Page - 1) * request.PerPage
		query = query.Offset(offset).Limit(request.PerPage)
	}

	if err := query.Find(&routesList).Error; err != nil {
		r.Log.WithError(err).Error("failed to find routes")
		return nil, 0, err
	}

	return routesList, total, nil
}

func (r *routeRepositoryImpl) FilterRoute(request *model.FindAllRouteRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// Filter by search
		if search := request.Search; search != "" {
			search = "%" + search + "%"
			tx = tx.Where("name LIKE ?", search)
		}

		return tx
	}
}

func (r *routeRepositoryImpl) Create(db *gorm.DB, route *entity.Route) error {
	return db.Create(route).Error
}

func (r *routeRepositoryImpl) Update(db *gorm.DB, route *entity.Route) error {
	return db.Model(&entity.Route{ID: route.ID}).Updates(route).Error
}

func (r *routeRepositoryImpl) Delete(db *gorm.DB, id int) error {
	return db.Delete(&entity.Route{}, id).Error
}

func (r *routeRepositoryImpl) CountByName(db *gorm.DB, name string) (int64, error) {
	var count int64
	if err := db.Model(&entity.Route{}).Where("name = ?", name).Count(&count).Error; err != nil {
		r.Log.WithError(err).Error("failed to count routes")
		return 0, err
	}
	return count, nil
}

func (r *routeRepositoryImpl) FindById(db *gorm.DB, id int) (*entity.Route, error) {
	var route entity.Route

	err := db.Preload("Sales").First(&route, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &route, nil
}
