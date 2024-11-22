package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type WishlistService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewWishlistService(repo repository.MainRepository, log *zap.Logger) WishlistService {
	return WishlistService{Repo: repo, Logger: log}
}

func (s *WishlistService) AddProductToWishlist(wishlistInput model.WishlistDTO) error {
	return s.Repo.WishlistRepository.Add(wishlistInput)
}

func (s *WishlistService) RemoveProductFromWishlist(userId string, wishlistId int) error {
	return s.Repo.WishlistRepository.Delete(userId, wishlistId)
}

func (s *WishlistService) GetWishlistByUserId(userId string, paginationInput model.Pagination) ([]model.Wishlist, model.Pagination, error) {
	var newWishlist []model.Wishlist
	if paginationInput.Page == 0 {
		paginationInput.Page = 1
	}
	if paginationInput.PerPage == 0 {
		paginationInput.PerPage = 5
	}
	wishlist, pagination, err := s.Repo.WishlistRepository.GetAll(userId, paginationInput)
	if err != nil {
		return nil, paginationInput, err
	}

	for _, item := range wishlist {
		product, err := s.Repo.ProductRepository.GetByID(item.ProductID)
		if err != nil {
			s.Logger.Error("Error getting product", zap.Error(err))
			continue
		}
		item.Product = product
		newWishlist = append(newWishlist, item)
	}
	return newWishlist, pagination, nil
}
