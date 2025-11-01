package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SalesRepository interface {
	Create(db *gorm.DB, sales *entity.Sales) error
	Update(db *gorm.DB, sales *entity.Sales) error
	Delete(db *gorm.DB, id int) error
	FindById(db *gorm.DB, id int) (*entity.Sales, error)
	FindAll(db *gorm.DB, request *model.FindAllSalesRequest) ([]entity.Sales, int64, error)
	CountByPhone(db *gorm.DB, phone string) (int64, error)
	ReplaceRoutes(db *gorm.DB, sales *entity.Sales, routeIDs []entity.Route) error
}

type salesRepositoryImpl struct {
	Log *logrus.Logger
}

func NewSalesRepository(log *logrus.Logger) SalesRepository {
	return &salesRepositoryImpl{
		Log: log,
	}
}

func (r *salesRepositoryImpl) Create(db *gorm.DB, sales *entity.Sales) error {
	return db.Create(sales).Error
}

func (r *salesRepositoryImpl) Update(db *gorm.DB, sales *entity.Sales) error {
	return db.Model(&entity.Sales{ID: sales.ID}).Updates(sales).Error
}

func (r *salesRepositoryImpl) Delete(db *gorm.DB, id int) error {
	return db.Delete(&entity.Sales{}, id).Error
}

func (r *salesRepositoryImpl) FindById(db *gorm.DB, id int) (*entity.Sales, error) {
	var sales entity.Sales

	err := db.First(&sales, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &sales, nil
}

func (r *salesRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllSalesRequest) ([]entity.Sales, int64, error) {
	var salesList []entity.Sales
	var total int64

	// Count total with filters
	countQuery := db.Model(&entity.Sales{}).Scopes(r.FilterSales(request))
	if err := countQuery.Count(&total).Error; err != nil {
		r.Log.WithError(err).Error("failed to count sales")
		return nil, 0, err
	}

	// Main query with filters, preload, and pagination
	query := db.Model(&entity.Sales{}).
		Scopes(r.FilterSales(request)).
		Preload("Routes").
		Order("name DESC")

	// Pagination
	if request.Page > 0 && request.PerPage > 0 {
		offset := (request.Page - 1) * request.PerPage
		query = query.Offset(offset).Limit(request.PerPage)
	}

	if err := query.Find(&salesList).Error; err != nil {
		r.Log.WithError(err).Error("failed to find sales")
		return nil, 0, err
	}

	return salesList, total, nil
}

func (r *salesRepositoryImpl) FilterSales(request *model.FindAllSalesRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		// Filter by search
		if search := request.Search; search != "" {
			search = "%" + search + "%"
			tx = tx.Where("name LIKE ? OR phone LIKE ?", search, search)
		}

		// Filter by multiple route IDs
		if len(request.RouteIDs) > 0 {
			tx = tx.Joins("JOIN sales_routes ON sales_routes.sales_id = sales.id").
				Where("sales_routes.route_id IN ?", request.RouteIDs).
				Distinct()
		}

		return tx
	}
}

func (r *salesRepositoryImpl) CountByPhone(db *gorm.DB, phone string) (int64, error) {
	var count int64
	err := db.Model(&entity.Sales{}).Where("phone = ?", phone).Count(&count).Error

	return count, err
}

func (r *salesRepositoryImpl) ReplaceRoutes(db *gorm.DB, sales *entity.Sales, routeIDs []entity.Route) error {

	err := db.Model(sales).Association("Routes").Replace(routeIDs)

	return err
}
