package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type MainRepository struct {
	UserRepository     UserRepository
	CategoryRepository CategoryRepository
	ProductRepository  ProductRepository
	VariantRepository  VariantRepository
}

func NewMainRepository(db *sql.DB, log *zap.Logger) MainRepository {
	return MainRepository{
		UserRepository:     NewUserRepository(db, log),
		CategoryRepository: NewCategoryRepository(db, log),
		ProductRepository:  NewProductRepository(db, log),
		VariantRepository:  NewVariantRepository(db, log),
	}
}
