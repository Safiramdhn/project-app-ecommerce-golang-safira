package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type MainRepository struct {
	AddressRepository        AddressRepository
	CategoryRepository       CategoryRepository
	OrderRepository          OrderRepository
	ProductRepository        ProductRepository
	RecommendationRepository RecommendationRepository
	UserRepository           UserRepository
	VariantRepository        VariantRepository
	WishlistRepository       WishlistRepository
	CartRepository           CartRepository
}

func NewMainRepository(db *sql.DB, log *zap.Logger) MainRepository {
	return MainRepository{
		AddressRepository:        NewAddressRepository(db, log),
		CategoryRepository:       NewCategoryRepository(db, log),
		OrderRepository:          NewOrderRepository(db, log),
		ProductRepository:        NewProductRepository(db, log),
		RecommendationRepository: NewRecommendationRepository(db, log),
		UserRepository:           NewUserRepository(db, log),
		VariantRepository:        NewVariantRepository(db, log),
		WishlistRepository:       NewWishlistRepository(db, log),
		CartRepository:           NewCartRepository(db, log),
	}
}
