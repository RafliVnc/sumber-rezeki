package usecase

import (
	"api/internal/entity"
	"api/internal/model"
	"api/internal/model/converter"
	"api/internal/repository"
	"api/internal/utils"
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	FindAll(ctx context.Context, request *model.FindAllUserRequest) ([]model.UserResponse, int64, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
}

type UserUseCaseImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository repository.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

func (s *UserUseCaseImpl) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//check email
	count, err := s.UserRepository.CountByEmail(tx, request.Email)
	if err != nil {
		s.Log.Warnf("Failed check email to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		s.Log.Warnf("Email already exists : %s", request.Email)
		errorMessage := fmt.Sprintf("Email %s already exists", request.Email)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//encript password
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Log.Warnf("Failed to generate bcrypt hash : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//set user
	user := &entity.User{
		Email:    request.Email,
		Password: string(password),
		Name:     request.Name,
	}

	if err := s.UserRepository.Create(tx, user); err != nil {
		s.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"email": request.Email,
			"name":  request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToUserResponse(user), nil

}

func (s *UserUseCaseImpl) FindAll(ctx context.Context, request *model.FindAllUserRequest) ([]model.UserResponse, int64, error) {
	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, 0, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	users, total, err := s.UserRepository.FindAll(s.DB.WithContext(ctx), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting users")
		return nil, 0, fiber.ErrInternalServerError
	}

	responses := make([]model.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *converter.ToUserResponse(&user)
	}

	return responses, total, nil
}

func (s *UserUseCaseImpl) Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//check email
	count, err := s.UserRepository.CheckEmailUniqueness(tx, request.Email, request.ID)
	if err != nil {
		s.Log.Warnf("Failed check email to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		s.Log.Warnf("Email already exists : %s", request.Email)
		errorMessage := fmt.Sprintf("Email %s already exists", request.Email)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//set user
	user := &entity.User{
		ID:    request.ID,
		Email: request.Email,
		Name:  request.Name,
	}

	if err := s.UserRepository.Update(tx, user); err != nil {
		s.Log.Warnf("Failed update user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"email": request.Email,
			"name":  request.Name,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToUserResponse(user), nil
}

func (s *UserUseCaseImpl) Delete(ctx context.Context, request *model.DeleteUserRequest) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//Check if user exists
	user, err := s.UserRepository.FindById(tx, request.ID)
	if err != nil {
		s.Log.WithError(err).Error("error getting user")
		return fiber.ErrInternalServerError
	}

	if user == nil {
		s.Log.Warnf("User not found : %d", request.ID)
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("user with id %d not found", request.ID))
	}

	if err := s.UserRepository.Delete(tx, request.ID); err != nil {
		s.Log.WithError(err).Error("error deleting user")
		return fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"id": request.ID,
		}).Warnf("Failed commit to database : %+v", err)
		return fiber.ErrInternalServerError
	}

	return nil
}
