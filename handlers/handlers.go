package handlers

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"go.uber.org/zap"
)

type Mainhandler struct {
	AddressHandler        AddressHandler
	CategoryHandler       CategoryHandler
	CartHandler           CartHandler
	ProductHandler        ProductHandler
	RecommendationHandler RecommendationHandler
	UserHandler           UserHandler
	WishlistHandler       WishlistHandler
}

func NewMainHandler(service service.MainService, log *zap.Logger, config util.Configuration) Mainhandler {
	return Mainhandler{
		AddressHandler:        NewAddressHandler(service, log),
		CategoryHandler:       NewCategoryHandler(service, log),
		CartHandler:           NewCartHandler(service, log),
		ProductHandler:        NewProductHandler(service, log),
		RecommendationHandler: NewRecommendationHandler(service, log),
		UserHandler:           NewUserHandler(service, log, config),
		WishlistHandler:       NewWishlistHandler(service, log),
	}
}
