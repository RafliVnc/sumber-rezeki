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
	Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error)
	Delete(ctx context.Context, request *model.DeleteUserRequest) error
	Verify(ctx context.Context, request *model.Auth) (*model.Auth, error)
	Current(ctx context.Context, jwtToken string) (*model.UserResponse, error)
}

type UserUseCaseImpl struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository repository.UserRepository
	TokenUtil      *utils.TokenUtil
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository repository.UserRepository, tokenUtil *utils.TokenUtil) *UserUseCaseImpl {
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
	if err != nil {
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
		errorMessage := fmt.Sprintf("Username %s already exists", request.Username)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	// TODO: CHECK PHONE

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

	//check username
	count, err := s.UserRepository.CheckUsernameUniqueness(tx, request.Username, request.ID)
	if err != nil {
		s.Log.Warnf("Failed check username to database : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	if count > 0 {
		s.Log.Warnf("Username already exists : %s", request.Username)
		errorMessage := fmt.Sprintf("Username %s already exists", request.Username)
		return nil, fiber.NewError(fiber.StatusBadRequest, errorMessage)
	}

	//set user
	user := &entity.User{
		ID:       request.ID,
		Username: request.Username,
		Name:     request.Name,
	}

	if err := s.UserRepository.Update(tx, user); err != nil {
		s.Log.Warnf("Failed update user to database : %+v", err)
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

func (s *UserUseCaseImpl) Login(ctx context.Context, request *model.LoginUserRequest) (*model.LoginUserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	//check request validation
	details, errorMessage, err := utils.ValidateStruct(request)
	if err != nil {
		s.Log.Warnf("Failed to validate request: %+v", details)
		return nil, model.NewErrorResponse(fiber.StatusBadRequest, errorMessage, details)
	}

	user, err := s.UserRepository.FindByUsername(tx, request.Username)
	if err != nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		s.Log.Warnf("Failed compare password : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	token, err := s.TokenUtil.CreateToken(ctx, &model.Auth{ID: user.ID})
	if err != nil {
		s.Log.Warnf("Failed to create token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return converter.UserToTokenResponse(user, token), nil
}

func (s *UserUseCaseImpl) Verify(ctx context.Context, request *model.Auth) (*model.Auth, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := s.Validate.Struct(request)
	if err != nil {
		s.Log.Warnf("Invalid request body : %+v", err)
		return nil, fiber.ErrBadRequest
	}

	user, err := s.UserRepository.FindById(tx, request.ID)
	if err != nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	if err := tx.Commit().Error; err != nil {
		s.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	return &model.Auth{ID: user.ID}, nil
}

func (s *UserUseCaseImpl) Current(ctx context.Context, jwtToken string) (*model.UserResponse, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	token, err := s.TokenUtil.ParseToken(ctx, jwtToken)
	if err != nil {
		s.Log.Warnf("Failed to parse token : %+v", err)
		return nil, fiber.ErrInternalServerError
	}

	user, err := s.UserRepository.FindById(tx, token.ID)
	if err != nil {
		s.Log.Warnf("Failed find user by username : %+v", err)
		return nil, fiber.ErrUnauthorized
	}

	return converter.ToUserResponse(user), nil
}
