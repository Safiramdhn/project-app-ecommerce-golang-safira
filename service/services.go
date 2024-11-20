package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type MainService struct {
	UserService     UserService
	CategoryService CategoryService
}

func NewMainService(repo repository.MainRepository, log *zap.Logger) MainService {
	return MainService{
		UserService:     NewUserService(repo, log),
		CategoryService: NewCategoryService(repo, log),
	}
}
