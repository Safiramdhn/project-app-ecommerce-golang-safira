package repository

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type UserRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{DB: db, Logger: logger}
}

func (repo *UserRepository) Create(userInput model.User) error {
	return nil
}

func (repo *UserRepository) Login(user model.User) (*model.User, error) {
	return nil, nil
}

func (repo *UserRepository) GetByID(id int) (*model.User, error) {
	return nil, nil
}

func (repo *UserRepository) Update(user *model.User) error {
	return nil
}

func (repo *UserRepository) Delete(id int) error {
	return nil
}
