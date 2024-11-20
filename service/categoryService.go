package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type CategoryService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewCategoryService(repo repository.MainRepository, logger *zap.Logger) CategoryService {
	return CategoryService{Repo: repo, Logger: logger}
}

func (s *CategoryService) GetAllCategory(pagination model.Pagination) ([]model.Category, model.Pagination, error) {
	if pagination.Page == 0 && pagination.PerPage == 0 {
		pagination.Page = 1
		pagination.PerPage = 5
	}
	return s.Repo.CategoryRepository.GetAll(pagination)
}
