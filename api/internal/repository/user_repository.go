package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(db *gorm.DB, entity *entity.User) error
	Update(db *gorm.DB, entity *entity.User) error
	CountByUsername(db *gorm.DB, username string) (int64, error)
	FindAll(db *gorm.DB, request *model.FindAllUserRequest) ([]entity.User, int64, error)
	CheckUsernameUniqueness(tx *gorm.DB, username string, id uuid.UUID) (int64, error)
	Delete(db *gorm.DB, id uuid.UUID) error
	FindById(db *gorm.DB, id uuid.UUID) (*entity.User, error)
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

func (r *UserRepositoryImpl) Update(db *gorm.DB, entity *entity.User) error {
	return db.Save(entity).Error
}

func (r *UserRepositoryImpl) CountByUsername(db *gorm.DB, username string) (int64, error) {
	var count int64
	err := db.Model(new(entity.User)).Where("username = ?", username).Count(&count).Error
	return count, err
}

func (r *UserRepositoryImpl) FindAll(db *gorm.DB, request *model.FindAllUserRequest) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64

	query := db.Model(new(entity.User)).Scopes(r.FilterUser(request))
	if request.Page > 0 && request.PerPage > 0 {
		query = query.Offset((request.Page - 1) * request.PerPage).Limit(request.PerPage)
	}

	if err := query.Count(&total).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) FilterUser(request *model.FindAllUserRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if username := request.Username; username != "" {
			username = "%" + username + "%"
			tx = tx.Where("username LIKE ?", username)
		}

		return tx
	}
}

func (r *UserRepositoryImpl) CheckUsernameUniqueness(tx *gorm.DB, username string, id uuid.UUID) (int64, error) {
	var count int64
	if err := tx.Model(new(entity.User)).Where("username = ? AND id != ?", username, id).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
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
