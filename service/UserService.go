package service

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewUserService(repo repository.MainRepository, log *zap.Logger) UserService {
	return UserService{Repo: repo, Logger: log}
}

func (s *UserService) CreateUser(userInput model.UserDTO) (string, error) {
	//encode password
	passwordHashed, err := helper.EncodePassword(userInput.Password)
	if err != nil {
		s.Logger.Error("error encode password", zap.Error(err), zap.String("Service", "User"), zap.String("Function", "CreateUser"))
		return "", err
	}
	newUserInput := model.User{
		ID:             uuid.NewString(),
		Name:           userInput.Name,
		PasswordHashed: passwordHashed,
		Email:          sql.NullString{String: userInput.EmailOrPhoneNumber, Valid: true},
		PhoneNumber:    sql.NullString{String: userInput.EmailOrPhoneNumber, Valid: true},
	}

	err = s.Repo.UserRepository.Create(newUserInput)
	if err != nil {
		s.Logger.Error("error creating user", zap.Error(err), zap.String("Service", "User"), zap.String("Function", "CreateUser"))
		return "", err
	}
	return newUserInput.ID, nil
}

func (s *UserService) Login(userInput model.UserDTO) (model.User, error) {
	user, err := s.Repo.UserRepository.Login(userInput)
	if err != nil {
		s.Logger.Error("error login user", zap.Error(err), zap.String("Service", "User"), zap.String("Function", "Login"))
		return model.User{}, err
	}
	if user.ID == "" {
		return model.User{}, nil
	}

	// compare password
	if user.PasswordHashed != "" {
		passwordValidation, err := helper.ComparePassword(user.PasswordHashed, userInput.Password)
		if !passwordValidation {
			s.Logger.Error("password validation failed", zap.Error(err), zap.String("Service", "User"), zap.String("Function", "Login"))
			return model.User{}, err
		}
	}
	return user, nil
}
