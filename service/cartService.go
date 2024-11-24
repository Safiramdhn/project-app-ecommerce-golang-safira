package service

import (
	"errors"

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

func (s *CartService) AddProductToCart(userID string, CartInput model.CartItemDTO) error {
	// Get the existing cart for the user
	existedCart, err := s.Repo.CartRepository.GetByUserID(userID)
	if err != nil {
		s.Logger.Error("Failed to get existing cart", zap.Error(err))
		return err
	}

	var cart model.Cart

	// If the cart doesn't exist, create a new one
	if existedCart.ID == 0 {
		s.Logger.Info("Creating new cart", zap.String("service", "cart"), zap.String("function", "AddProductToCart"))
		cart, err = s.Repo.CartRepository.Create(userID)
		if err != nil {
			s.Logger.Error("Failed to create new cart", zap.Error(err))
			return err
		}
	} else {
		// If the cart already exists, use the existing cart
		s.Logger.Info("Cart already exists", zap.String("service", "cart"), zap.String("function", "AddProductToCart"))
		cart = existedCart
	}

	product, err := s.Repo.ProductRepository.GetByID(CartInput.ProductID)
	if err != nil {
		s.Logger.Error("error get product by id", zap.Error(err))
		return err
	}

	weekly, err := s.Repo.ProductRepository.GetPromoProduct(product.ID)
	if err != nil {
		s.Logger.Error("error get weekly promo by product id", zap.Error(err))
		return err
	}
	if weekly.ID != 0 {
		product.PriceAfterDiscount = helper.CalculateDiscountPrice(product.PriceAfterDiscount, weekly.PromoDiscount)
	}
	if CartInput.Amount == 0 {
		CartInput.Amount = 1
	}
	cartItem := model.CartItem{
		ProductID: CartInput.ProductID,
		Amount:    CartInput.Amount,
		SubTotal:  product.PriceAfterDiscount * float64(CartInput.Amount),
		Product:   product,
		CartID:    cart.ID,
	}
	additionalPrice, err := s.AddItemToCart(cartItem, CartInput.Variant)
	if err != nil {
		s.Logger.Error("error add item to second cart", zap.Error(err))
		return err
	}

	cart.TotalAmount += cartItem.Amount
	cart.TotalPrice += (additionalPrice + product.PriceAfterDiscount) * float64(CartInput.Amount)
	err = s.UpdateCart(cart)
	if err != nil {
		s.Logger.Error("error update second cart", zap.Error(err))
		return err
	}
	return nil
}

func (s *CartService) GetCartByUserID(userID string) (model.Cart, error) {
	cart, err := s.Repo.CartRepository.GetByUserID(userID)
	if err != nil {
		return model.Cart{}, err
	}
	cartItems, err := s.Repo.CartRepository.GetItems(cart.ID)
	if err != nil {
		return model.Cart{}, err
	}

	var newCartItem []model.CartItem
	for _, item := range cartItems {
		product, err := s.Repo.ProductRepository.GetByID(item.ProductID)
		if err != nil {
			s.Logger.Error("error get product by id", zap.Error(err))
			return model.Cart{}, err
		}
		item.Product = product
		NewVariants := []model.CarttemVariant{}
		if product.HasVariant {
			for _, variantItem := range item.ItemVariant {
				variant, err := s.Repo.VariantRepository.GetVariantByID(int(variantItem.VariantID.Int64))
				if err != nil {
					s.Logger.Error("error get variant by id", zap.Error(err))
					return model.Cart{}, err
				}
				option, err := s.Repo.VariantRepository.GetVariantOptionByID(int(variantItem.OptionID.Int64))
				if err != nil {
					s.Logger.Error("error get variant option by id", zap.Error(err))
					return model.Cart{}, err
				}

				variantItem.Variant = variant
				variantItem.Option = option
				NewVariants = append(NewVariants, variantItem)
			}
			item.ItemVariant = NewVariants
		}
		newCartItem = append(newCartItem, item)
	}
	cart.Items = newCartItem
	return cart, nil
}

func (s *CartService) AddItemToCart(itemInput model.CartItem, itemVariantInput []model.CartItemVariantDTO) (float64, error) {
	var AdditionalPrice = 0.0
	item, err := s.Repo.CartRepository.AddItem(itemInput)
	if err != nil {
		s.Logger.Error("error add item to second cart", zap.Error(err))
		return AdditionalPrice, err
	}
	if item.Product.HasVariant {
		for _, v := range itemVariantInput {
			err := s.AddVariantItem(item.ID, v)
			if err != nil {
				s.Logger.Error("error add variant item to second cart", zap.Error(err))
				return AdditionalPrice, err
			}

			for _, pv := range item.Product.Variant {
				if v.VariantID == pv.ID {
					for _, option := range pv.VariantOption {
						if v.VariantOptionID == option.ID {
							AdditionalPrice += option.AdditionalPrice
							break
						}
					}
				}
			}
		}

	}
	return AdditionalPrice, nil
}

func (s *CartService) AddVariantItem(cartItemID int, variantInput model.CartItemVariantDTO) error {
	varianOption, err := s.Repo.VariantRepository.GetVariantOptionByID(variantInput.VariantOptionID)
	if err != nil {
		s.Logger.Error("error get variant option", zap.Error(err))
		return err
	}
	variantInput.AdditionalPrice = varianOption.AdditionalPrice
	return s.Repo.CartRepository.AddItemVariant(cartItemID, variantInput)
}

func (s *CartService) UpdateCart(cartInput model.Cart) error {
	return s.Repo.CartRepository.Update(cartInput)
}

func (s *CartService) UpdateItemInCart(userId string, itemInput model.CartItem) error {
	cart, err := s.Repo.CartRepository.GetByID(itemInput.CartID)
	if err != nil {
		s.Logger.Error("error get cart by id", zap.Error(err))
		return err
	}
	if cart.ID == 0 {
		s.Logger.Error("cart not found", zap.Error(err))
		return errors.New("cart not found")
	}
	if cart.UserID != userId {
		s.Logger.Error("UserID mismatch for cart", zap.Error(err))
		return errors.New("UserID mismatch for cart")
	}

	product, err := s.Repo.ProductRepository.GetByID(itemInput.ProductID)
	if err != nil {
		s.Logger.Error("error get product by id", zap.Error(err))
		return err
	}

	weekly, err := s.Repo.ProductRepository.GetPromoProduct(product.ID)
	if err != nil {
		s.Logger.Error("error get weekly promo by product id", zap.Error(err))
		return err
	}
	if weekly.ID != 0 {
		product.PriceAfterDiscount = helper.CalculateDiscountPrice(product.PriceAfterDiscount, weekly.PromoDiscount)
	}

	var additional_costs float64
	if product.HasVariant {
		itemVariant, err := s.Repo.CartRepository.GetItemVariants(itemInput.ID)
		if err != nil {
			s.Logger.Error("error get item variants", zap.Error(err))
			return err
		}

		if len(itemVariant) > 0 {
			for _, variant := range itemVariant {
				additional_costs += variant.AdditionalPrice
			}
		}
	}

	itemInput.SubTotal = (product.PriceAfterDiscount + additional_costs) * float64(itemInput.Amount)
	err = s.Repo.CartRepository.UpdateItem(itemInput)
	if err != nil {
		s.Logger.Error("error update item in cart", zap.Error(err))
		return err
	}
	err = s.Repo.CartRepository.RecalculateTotal(cart.ID)
	if err != nil {
		s.Logger.Error("error recalculate total amount in cart", zap.Error(err))
		return err
	}
	return nil
}

func (s *CartService) DeleteProductInCart(cartItemID int) error {
	cartItem, err := s.Repo.CartRepository.GetItemByID(cartItemID)
	if err != nil {
		s.Logger.Error("error get cart item by id", zap.Error(err))
		return err
	}
	if cartItem.CartID == 0 {
		s.Logger.Error("cart not found", zap.String("service", "cart"), zap.String("function", "DeleteProductInCart"))
		return errors.New("cart not found")
	}

	err = s.Repo.CartRepository.DeleteItem(cartItemID)
	if err != nil {
		s.Logger.Error("error delete item from cart", zap.Error(err))
		return err
	}
	err = s.Repo.CartRepository.RecalculateTotal(cartItem.CartID)
	if err != nil {
		s.Logger.Error("error recalculate total amount in cart", zap.Error(err))
		return err
	}
	return nil
}
