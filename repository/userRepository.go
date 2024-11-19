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
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err))
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO users (id,name, email, phone_number, password)
			VALUES ($1, $2 CASE WHEN $3 ~ '^[^@]+@[^@]+\.[^@]+$' THEN $3 ELSE NULL END, 
			CASE WHEN $4 ~ '^[0-9]{10,15}$' THEN $4 ELSE NULL END, 
			$3, $4) RETURNING id;`

	_, err = repo.DB.Exec(sqlStatement, userInput.ID, userInput.Name, userInput.Email, userInput.PhoneNumber, userInput.PasswordHashed)
	if err != nil {
		repo.Logger.Error("Error creating user", zap.Error(err))
		return err
	}
	return nil
}

func (repo *UserRepository) Login(userLogin model.UserDTO) (*model.User, error) {
	var user model.User
	sqlStatement := `SELECT id, password FROM users WHERE (email = $1 OR phone_number = $1)`

	err := repo.DB.QueryRow(sqlStatement, userLogin.EmailOrPhoneNumber).Scan(&user.ID, &user.PasswordHashed)
	if err == sql.ErrNoRows {
		repo.Logger.Error("User not found", zap.Error(err))
		return &user, nil
	} else if err != nil {
		repo.Logger.Error("Error retrieving user", zap.Error(err))
		return nil, err
	}

	return &user, nil
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
