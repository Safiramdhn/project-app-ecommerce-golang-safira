package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type MainRepository struct {
	UserRepository           UserRepository
	CategoryRepository       CategoryRepository
	ProductRepository        ProductRepository
	VariantRepository        VariantRepository
	RecommendationRepository RecommendationRepository
	WishlistRepository       WishlistRepository
}

func NewMainRepository(db *sql.DB, log *zap.Logger) MainRepository {
	return MainRepository{
		UserRepository:           NewUserRepository(db, log),
		CategoryRepository:       NewCategoryRepository(db, log),
		ProductRepository:        NewProductRepository(db, log),
		VariantRepository:        NewVariantRepository(db, log),
		RecommendationRepository: NewRecommendationRepository(db, log),
		WishlistRepository:       NewWishlistRepository(db, log),
	}
}
