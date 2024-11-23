package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type CartService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewCartService(repo repository.MainRepository, logger *zap.Logger) CartService {
	return CartService{Repo: repo, Logger: logger}
}

func (s CartService) AddCart(userID string, cartInput model.CartDTO) error {
	newCartInput := model.Cart{
		ProductID:        cartInput.ProductID,
		VariantIDs:       helper.ConvertToNullInt64Slice(cartInput.VariantIDs),
		VariantOptionIDs: helper.ConvertToNullInt64Slice(cartInput.VariantOptionIDs),
		Amount:           cartInput.Amount,
	}
	additionalPrice := 0.0

	product, err := s.Repo.ProductRepository.GetByID(cartInput.ProductID)
	if err != nil {
		s.Logger.Error("error get product by id", zap.Error(err))
		return err
	}

	if product.HasVariant && len(cartInput.VariantOptionIDs) > 0 {
		for _, v := range cartInput.VariantOptionIDs {
			variantOption, err := s.Repo.VariantRepository.GetVariantOptionByID(v)
			if err != nil {
				s.Logger.Error("error get variant by id", zap.Error(err))
				return err
			}
			if variantOption.AdditionalPrice != 0.0 {
				additionalPrice += variantOption.AdditionalPrice
			}
		}
	}

	var promoDiscount = 0.0
	weeklyPromo, err := s.Repo.ProductRepository.GetPromoProduct(product.ID)
	if err != nil {
		s.Logger.Error("error get weekly promo by product id", zap.Error(err))
		return err
	}

	if weeklyPromo.ID > 0 {
		promoDiscount = weeklyPromo.PromoDiscount
	}

	newCartInput.TotalPrice = helper.CalculateCartPrice(product.Price, additionalPrice, product.Discount, promoDiscount, newCartInput.Amount)

	return s.Repo.CartRepository.Create(userID, newCartInput)
}

func (s CartService) GetCartByUserID(userID string) ([]model.Cart, error) {
	cart, err := s.Repo.CartRepository.GetAll(userID)
	if err != nil {
		s.Logger.Error("error get cart by user id", zap.Error(err))
		return nil, err
	}

	var newCart []model.Cart
	for _, v := range cart {
		product, err := s.Repo.ProductRepository.GetByID(v.ProductID)
		if err != nil {
			s.Logger.Error("error get product by id", zap.Error(err))
			return nil, err
		}
		v.Product = product
		v.Product.PriceAfterDiscount = helper.CalculateDiscountPrice(v.Product.Price, v.Product.Discount)

		if v.Product.HasVariant {
			variants := []model.Variant{}
			for _, variantId := range v.VariantIDs {
				id := int(variantId.Int64)
				variant, err := s.Repo.VariantRepository.GetVariantByID(id)
				if err != nil {
					s.Logger.Error("error get variant by product id", zap.Error(err))
					return nil, err
				}
				variants = append(variants, variant)
			}
			v.Variants = variants

			variantOptions := []model.VariantOption{}
			for _, variantId := range v.VariantOptionIDs {
				id := int(variantId.Int64)
				variantOption, err := s.Repo.VariantRepository.GetVariantOptionByID(id)
				if err != nil {
					s.Logger.Error("error get variant option by id", zap.Error(err))
					return nil, err
				}
				variantOptions = append(variantOptions, variantOption)
			}
			v.VariantOptions = variantOptions
		}
		newCart = append(newCart, v)
	}
	return newCart, nil
}

func (s CartService) UpdateCart(userID string, cartID int, cartInput model.CartDTO) error {
	newCartInput := model.Cart{
		VariantOptionIDs: helper.ConvertToNullInt64Slice(cartInput.VariantOptionIDs),
		Amount:           cartInput.Amount,
	}
	additionalPrice := 0.0

	product, err := s.Repo.ProductRepository.GetByID(cartInput.ProductID)
	if err != nil {
		s.Logger.Error("error get product by id", zap.Error(err))
		return err
	}

	if product.HasVariant && len(cartInput.VariantOptionIDs) > 0 {
		for _, v := range cartInput.VariantOptionIDs {
			variantOption, err := s.Repo.VariantRepository.GetVariantOptionByID(v)
			if err != nil {
				s.Logger.Error("error get variant by id", zap.Error(err))
				return err
			}
			if variantOption.AdditionalPrice != 0.0 {
				additionalPrice += variantOption.AdditionalPrice
			}
		}
	}
	var promoDiscount = 0.0
	weeklyPromo, err := s.Repo.ProductRepository.GetPromoProduct(product.ID)
	if err != nil {
		s.Logger.Error("error get weekly promo by product id", zap.Error(err))
		return err
	}

	if weeklyPromo.ID > 0 {
		promoDiscount = weeklyPromo.PromoDiscount
	}
	newCartInput.TotalPrice = helper.CalculateCartPrice(product.Price, additionalPrice, product.Discount, promoDiscount, newCartInput.Amount)

	return s.Repo.CartRepository.Update(cartID, userID, newCartInput)
}

func (s *CartService) DeleteProductInCart(cartID int, userID string) error {
	return s.Repo.CartRepository.Delete(cartID, userID)
}

func (s *CartService) GetTotalCart() (int, float64, error) {
	return s.Repo.CartRepository.CountItemCart()
}
