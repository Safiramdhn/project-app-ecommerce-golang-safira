package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type MainService struct {
	AddressService        AddressService
	CategoryService       CategoryService
	ProductService        ProductService
	RecommendationService RecommendationService
	UserService           UserService
	WishlistService       WishlistService
	CartService           CartService
	OrderService          OrderService
}

func NewMainService(repo repository.MainRepository, log *zap.Logger) MainService {
	return MainService{
		AddressService:        NewAddressService(repo, log),
		CategoryService:       NewCategoryService(repo, log),
		ProductService:        NewProductService(repo, log),
		RecommendationService: NewRecommendationService(repo, log),
		UserService:           NewUserService(repo, log),
		WishlistService:       NewWishlistService(repo, log),
		CartService:           NewCartService(repo, log),
		OrderService:          NewOrderService(repo, log),
	}
}
