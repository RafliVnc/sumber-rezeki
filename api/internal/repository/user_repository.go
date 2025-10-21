package repository

import (
	"api/internal/entity"
	"api/internal/entity/enum"
	"api/internal/model"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(db *gorm.DB, request *model.FindAllUserRequest) ([]entity.User, int64, error)
	Create(db *gorm.DB, entity *entity.User) error
	Update(db *gorm.DB, id uuid.UUID, updates any) error
	Delete(db *gorm.DB, id uuid.UUID) error
	FindById(db *gorm.DB, id uuid.UUID) (*entity.User, error)
	CountByUsername(db *gorm.DB, username string) (int64, error)
	CountByPhone(db *gorm.DB, phone string) (int64, error)
	FindByUsername(db *gorm.DB, username string) (*entity.User, error)
}

type UserRepositoryImpl struct {
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		Log: log,
	}
}

func (r *UserRepositoryImpl) Create(db *gorm.DB, entity *entity.User) error {
	return db.Create(entity).Error
}

func (r *UserRepositoryImpl) Update(db *gorm.DB, id uuid.UUID, updates any) error {
	return db.Model(&entity.User{}).Where("id = ?", id).Updates(updates).Error
}

func (r *UserRepositoryImpl) CountByUsername(db *gorm.DB, username string) (int64, error) {
	var count int64

	// Cek with soft delete
	err := db.Unscoped().Model(new(entity.User)).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) CountByPhone(db *gorm.DB, phone string) (int64, error) {
	var count int64

	// Cek with soft delete
	err := db.Unscoped().Model(new(entity.User)).Where("phone = ?", phone).Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllUserRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	countQuery := db.Model(new(entity.User)).Scopes(r.FilterUser(request))
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := countQuery
	if request.Page > 0 && request.PerPage > 0 {
		query = query.Offset((request.Page - 1) * request.PerPage).Limit(request.PerPage)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) FilterUser(request *model.FindAllUserRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if search := request.Name; search != "" && request.Name == request.Username && request.Username == request.Phone {
			search = "%" + search + "%"
			tx = tx.Where("name LIKE ? OR username LIKE ? OR phone LIKE ?", search, search, search)
		}

		if len(request.Roles) > 0 {
			tx = tx.Where("role IN ?", request.Roles)
		}

		// exclude SUPER_ADMIN
		tx = tx.Where("role != ?", enum.SUPER_ADMIN)

		return tx
	}
}

func (r *UserRepositoryImpl) Delete(db *gorm.DB, id uuid.UUID) error {
	return db.Delete(new(entity.User), id).Error
}

func (r *UserRepositoryImpl) FindById(db *gorm.DB, id uuid.UUID) (*entity.User, error) {
	user := &entity.User{}
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	user := &entity.User{}
	if err := db.First(&user, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
