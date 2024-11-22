package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type ProductService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewProductService(repo repository.MainRepository, logger *zap.Logger) ProductService {
	return ProductService{Repo: repo, Logger: logger}
}

func (s ProductService) GetAllProduct(productFilter model.ProductDTO, pagination model.Pagination) ([]model.Product, model.Pagination, error) {
	if pagination.Page == 0 {
		pagination.Page = 1
	}

	if pagination.PerPage == 0 {
		pagination.PerPage = 5
	}

	return s.Repo.ProductRepository.GetAll(productFilter, pagination)
}

func (s ProductService) GetProductByID(id int) (*model.Product, error) {
	product, err := s.Repo.ProductRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// get variant
	if product.HasVariant {
		variant, err := s.Repo.VariantRepository.GetByProductId(product.ID)
		if err != nil {
			s.Logger.Error("Error retrieving variant", zap.Error(err), zap.String("Service", "Product"), zap.String("Function", "GetProductByID"))
			return nil, err
		}

		product.Variant = append(product.Variant, variant...)
		return &product, nil
	}

	return &product, nil
}

func (s ProductService) GetPromoWeekly(paginationInput model.Pagination) ([]model.WeeklyPromo, model.Pagination, error) {
	if paginationInput.Page == 0 {
		paginationInput.Page = 1
	}

	if paginationInput.PerPage == 0 {
		paginationInput.PerPage = 5
	}

	weeklyPromo, pagination, err := s.Repo.ProductRepository.GetWeeklyPromo(paginationInput)
	if err != nil {
		return nil, paginationInput, err
	}

	var newWeeklyPromos []model.WeeklyPromo
	for _, item := range weeklyPromo {
		product, err := s.Repo.ProductRepository.GetByID(item.ProductID)
		if err != nil {
			s.Logger.Error("Error retrieving product", zap.Error(err), zap.String("Service", "Product"), zap.String("Function", "GetPromoWeekly"))
			return nil, paginationInput, err
		}
		item.Product = product
		item.PromoPrice = helper.CalculateDiscountPrice(product.PriceAfterDiscount, item.PromoDiscout)
		newWeeklyPromos = append(newWeeklyPromos, item)
		pagination.CountData++
	}
	return newWeeklyPromos, pagination, nil
}
