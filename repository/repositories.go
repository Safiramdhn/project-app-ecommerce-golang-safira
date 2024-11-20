package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type MainRepository struct {
	UserRepository UserRepository
}

func NewMainRepository(db *sql.DB, log *zap.Logger) MainRepository {
	return MainRepository{UserRepository: NewUserRepository(db, log)}
}
