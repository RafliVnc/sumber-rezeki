package repository

import (
	"api/internal/entity"
	"api/internal/model"
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(db *gorm.DB, entity *entity.User) error
	Update(db *gorm.DB, entity *entity.User) error
	CountByEmail(db *gorm.DB, email string) (int64, error)
	FindAll(db *gorm.DB, request *model.FindAllUserRequest) ([]entity.User, int64, error)
	CheckEmailUniqueness(tx *gorm.DB, email string, id int) (int64, error)
	Delete(db *gorm.DB, id int) error
	FindById(db *gorm.DB, id int) (*entity.User, error)
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

func (r *UserRepositoryImpl) CountByEmail(db *gorm.DB, email string) (int64, error) {
	var count int64
	err := db.Model(new(entity.User)).Where("email = ?", email).Count(&count).Error
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

		if email := request.Email; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("email LIKE ?", email)
		}

		return tx
	}
}

func (r *UserRepositoryImpl) CheckEmailUniqueness(tx *gorm.DB, email string, id int) (int64, error) {
	var count int64
	if err := tx.Model(new(entity.User)).Where("email = ? AND id != ?", email, id).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *UserRepositoryImpl) Delete(db *gorm.DB, id int) error {
	return db.Delete(new(entity.User), id).Error
}

func (r *UserRepositoryImpl) FindById(db *gorm.DB, id int) (*entity.User, error) {
	user := &entity.User{}
	if err := db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
