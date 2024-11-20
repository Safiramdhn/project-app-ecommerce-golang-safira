package handlers

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"go.uber.org/zap"
)

type Mainhandler struct {
	UserHandler     UserHandler
	CategoryHandler CategoryHandler
	ProductHandler  ProductHandler
}

func NewMainHandler(service service.MainService, log *zap.Logger, config util.Configuration) Mainhandler {
	return Mainhandler{
		UserHandler:     NewUserHandler(service, log, config),
		CategoryHandler: NewCategoryHandler(service, log, config),
		ProductHandler:  NewProductHandler(service, log, config),
	}
}
