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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error)
	FindAll(ctx context.Context, request *model.FindAllUserRequest) ([]model.UserResponse, int64, error)
	Update(ctx context.Context, request *model.UpdateUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, string, error)
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
	Verify(ctx context.Context, request *model.Auth) (*model.Auth, error)
	Current(ctx context.Context, id uuid.UUID) (*model.UserResponse, error)
}

type UserUseCaseImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository repository.UserRepository
	TokenUtil      *utils.TokenUtil
}

func NewUserUseCase(
	db *gorm.DB,
	logger *logrus.Logger,
	validate *validator.Validate,
	userRepository repository.UserRepository,
	tokenUtil *utils.TokenUtil,
) UserUseCase {
	return &UserUseCaseImpl{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
		TokenUtil:      tokenUtil,
	}
}

func (s *UserUseCaseImpl) Create(ctx context.Context, request *model.RegisterUserRequest) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil || details != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	//check username
	count, err := s.UserRepository.CountByUsername(tx, request.Username)
	if err != nil {
		s.Log.Warnf("Failed check username to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		s.Log.Warnf("Username already exists : %s", request.Username)
		errorMessage := fmt.Sprintf("Username %s sudah pernah digunakan", request.Username)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//check phone uniqueness
	phoneCount, err := s.UserRepository.CountByPhone(tx, request.Phone)
	if err != nil {
		s.Log.Warnf("Failed check phone to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if phoneCount > 0 {
		s.Log.Warnf("Phone already exists : %s", request.Phone)
		errorMessage := fmt.Sprintf("Nomor HP %s sudah digunakan", request.Phone)
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
		ID:       uuid.New(),
		Username: request.Username,
		Password: string(password),
		Name:     request.Name,
		Role:     request.Role,
		Phone:    request.Phone,
	}

	if err := s.UserRepository.Create(tx, user); err != nil {
		s.Log.Warnf("Failed create user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"username": request.Username,
			"name":     request.Name,
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

	//get users
	users, total, err := s.UserRepository.FindAll(s.DB.WithContext(ctx), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting users")
		return nil, 0, fiber.ErrInternalServerError
	}

	// convert to arry response
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

	// check user
	user, err := s.UserRepository.FindById(tx, request.ID)
	if err != nil {
		s.Log.Warnf("Failed find user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if user == nil {
		s.Log.Warnf("User not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "Pengguna tidak ditemukan")
	}

	//check username uniqueness
	count, err := s.UserRepository.CountByUsername(tx, request.Username)
	if err != nil {
		s.Log.Warnf("Failed check username to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 && user.Username != request.Username {
		s.Log.Warnf("Username already exists : %s", request.Username)
		errorMessage := fmt.Sprintf("Username %s sudah pernah digunakan", request.Username)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//check phone uniqueness
	phoneCount, err := s.UserRepository.CountByPhone(tx, request.Phone)
	if err != nil {
		s.Log.Warnf("Failed check phone to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if phoneCount > 0 && user.Phone != request.Phone {
		s.Log.Warnf("Phone already exists : %s", request.Phone)
		errorMessage := fmt.Sprintf("Nomor HP %s sudah digunakan", request.Phone)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//set user
	updateUser := &entity.User{
		ID:       request.ID,
		Name:     request.Name,
		Username: request.Username,
		Phone:    request.Phone,
		Role:     request.Role,
	}

	if err := s.UserRepository.Update(tx, user.ID, updateUser); err != nil {
		s.Log.Warnf("Failed update user to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		s.Log.WithFields(logrus.Fields{
			"username": request.Username,
			"name":     request.Name,
			"phone":    request.Phone,
			"role":     request.Role,
		}).Warnf("Failed commit to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.ToUserResponse(updateUser), nil
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
		return fiber.NewError(fiber.StatusNotFound, "Pengguna tidak ditemukan")
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

func (s *UserUseCaseImpl) Login(ctx context.Context, request *model.LoginUserRequest) (*model.UserResponse, string, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, "", model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	user, err := s.UserRepository.FindByUsername(tx, request.Username)
	if err != nil || user == nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Username atau kata sandi tidak valid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		s.Log.Warnf("Failed compare password : %+v", err)
		return nil, "", fiber.NewError(fiber.StatusUnauthorized, "Username atau kata sandi tidak valid")
	}

	token, err := s.TokenUtil.CreateToken(ctx, &model.Auth{ID: user.ID})
	if err != nil {
		s.Log.Warnf("Failed to create token : %+v", err)
		return nil, "", fiber.ErrInternalServerError
	}

	return converter.ToUserResponse(user), token, nil
}

func (s *UserUseCaseImpl) Verify(ctx context.Context, request *model.Auth) (*model.Auth, error) {
	err := s.Validate.Struct(request)
	if err != nil {
		s.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	// Chek user
	user, err := s.UserRepository.FindById(s.DB.WithContext(ctx), request.ID)
	if err != nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if user == nil {
		s.Log.Warnf("User not found : %d", request.ID)
		return nil, fiber.NewError(fiber.StatusNotFound, "pengguna tidak ditemukan")
	}

	return &model.Auth{ID: user.ID}, nil
}

func (s *UserUseCaseImpl) Current(ctx context.Context, id uuid.UUID) (*model.UserResponse, error) {

	user, err := s.UserRepository.FindById(s.DB.WithContext(ctx), id)
	if err != nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if user == nil {
		s.Log.Warnf("User not found : %d", id)
		return nil, fiber.NewError(fiber.StatusNotFound, "pengguna tidak ditemukan")
	}

	return converter.ToUserResponse(user), nil
}
