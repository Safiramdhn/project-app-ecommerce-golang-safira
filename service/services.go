package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type MainService struct {
	AddressService        AddressService
	CartService           CartService
	CategoryService       CategoryService
	ProductService        ProductService
	RecommendationService RecommendationService
	UserService           UserService
	WishlistService       WishlistService
}

func NewMainService(repo repository.MainRepository, log *zap.Logger) MainService {
	return MainService{
		AddressService:        NewAddressService(repo, log),
		CategoryService:       NewCategoryService(repo, log),
		CartService:           NewCartService(repo, log),
		ProductService:        NewProductService(repo, log),
		RecommendationService: NewRecommendationService(repo, log),
		UserService:           NewUserService(repo, log),
		WishlistService:       NewWishlistService(repo, log),
	}
}
