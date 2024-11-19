package service

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"github.com/google/uuid"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) CreateUser(userInput model.UserDTO) (string, error) {
	//encode password
	passwordHashed, err := helper.EncodePassword(userInput.Password)
	if err != nil {
		return "", err
	}
	newUserInput := model.User{
		ID:             uuid.NewString(),
		Name:           userInput.Name,
		PasswordHashed: passwordHashed,
		Email:          sql.NullString{String: userInput.EmailOrPhoneNumber, Valid: true},
		PhoneNumber:    sql.NullString{String: userInput.EmailOrPhoneNumber, Valid: true},
	}

	err = s.UserRepo.Create(newUserInput)
	if err != nil {
		return "", err
	}
	return newUserInput.ID, nil
}

func (s *UserService) Login(userInput model.UserDTO) (*model.User, error) {
	user, err := s.UserRepo.Login(userInput)
	if err != nil {
		return user, err
	}

	// compare password
	if user.PasswordHashed != "" {
		passwordValidation := helper.ComparePassword(user.PasswordHashed, userInput.Password)
		if !passwordValidation {
			return nil, err
		}
	}
	return user, nil
}
