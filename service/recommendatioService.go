package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type RecommendationService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewRecommendationService(repo repository.MainRepository, logger *zap.Logger) RecommendationService {
	return RecommendationService{Repo: repo, Logger: logger}
}

func (s *RecommendationService) GetProductRecommendations(recommedFilter model.RecommendationDTO, pagination model.Pagination) ([]model.Recommendation, model.Pagination, error) {
	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PerPage == 0 {
		pagination.PerPage = 5
	}

	return s.Repo.RecommendationRepository.GetRecommendations(recommedFilter, pagination)
}
