package service

import (
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/repository"
	"go.uber.org/zap"
)

type OrderService struct {
	Repo   repository.MainRepository
	Logger *zap.Logger
}

func NewOrderService(repo repository.MainRepository, logger *zap.Logger) OrderService {
	return OrderService{Repo: repo, Logger: logger}
}

func (s *OrderService) CreateOrder(userID string, orderInput model.OrderDTO) error {
	cart, err := s.Repo.CartRepository.GetByID(orderInput.CartID)
	if err != nil {
		s.Logger.Error("error get cart by id", zap.Error(err))
		return err
	}

	if cart.UserID != userID {
		s.Logger.Error("user not match")
		return err
	}

	newOrderInput := model.Order{
		UserID:        userID,
		CartID:        orderInput.CartID,
		TotalPrice:    cart.TotalPrice + orderInput.ShippingCost,
		TotalAmount:   cart.TotalAmount,
		AddressID:     orderInput.AddressID,
		ShippingType:  orderInput.ShippingType,
		ShippingCost:  orderInput.ShippingCost,
		PaymentMethod: orderInput.PaymentMethod,
	}

	order, err := s.Repo.OrderRepository.Create(newOrderInput)
	if err != nil {
		s.Logger.Error("error create order", zap.Error(err))
		return err
	}

	err = s.AddItemOrder(order)
	if err != nil {
		s.Logger.Error("error add item order", zap.Error(err))
		s.UpdateOrderStatus(order.ID, cart.ID, "failed")
		return err
	}
	s.UpdateOrderStatus(order.ID, cart.ID, "success")
	return nil
}

func (s *OrderService) AddItemOrder(order model.Order) error {
	cartItem, err := s.Repo.CartRepository.GetItems(order.CartID)
	if err != nil {
		s.Logger.Error("error get cart item by cart id", zap.Error(err))
		return err
	}
	for _, item := range cartItem {
		orderItemInput := model.OrderItem{
			OrderID:    order.ID,
			ProductID:  item.ProductID,
			Amount:     item.Amount,
			SubTotal:   item.SubTotal,
			CartItemID: item.ID,
		}
		orderItem, err := s.Repo.OrderRepository.AddOrderItem(orderItemInput)
		if err != nil {
			s.Logger.Error("error create order item", zap.Error(err))
			return err
		}

		err = s.AddVariantItem(orderItem)
		if err != nil {
			s.Logger.Error("error add variant item to order item", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *OrderService) AddVariantItem(orderItem model.OrderItem) error {
	variantItems, err := s.Repo.CartRepository.GetItemVariants(orderItem.CartItemID)
	if err != nil {
		s.Logger.Error("error get item variant by cart item id", zap.Error(err))
		return err
	}

	for _, variantItem := range variantItems {
		orderVariantItemInput := model.OrderItemVariant{
			OrderItemID: orderItem.ID,
			VariantID:   variantItem.VariantID,
			OptionID:    variantItem.OptionID,
		}
		err = s.Repo.OrderRepository.AddOrderItemVariant(orderVariantItemInput)
		if err != nil {
			s.Logger.Error("error create order item variant", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *OrderService) UpdateOrderStatus(orderID, cartID int, orderStatus string) error {
	err := s.Repo.OrderRepository.UpdateOrderStatus(orderID, orderStatus)
	if err != nil {
		s.Logger.Error("error update order status", zap.Error(err))
		return err
	}

	if orderStatus == "success" {
		err = s.Repo.CartRepository.UpdateCartStatus(cartID)
		if err != nil {
			s.Logger.Error("error delete cart after success order", zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *OrderService) GetOrderByID(orderID int) (*model.Order, error) {
	order, err := s.Repo.OrderRepository.GetByID(orderID)
	if err != nil {
		s.Logger.Error("error get order by id", zap.Error(err))
		return nil, err
	}

	address, err := s.Repo.AddressRepository.GetByID(order.AddressID)
	if err != nil {
		s.Logger.Error("error get address by id", zap.Error(err))
		return nil, err
	}
	order.Address = *address
	orderItems, err := s.GetOrderItems(order.ID)
	if err != nil {
		s.Logger.Error("error get order item by order id", zap.Error(err))
		return nil, err
	}
	order.OrderItems = orderItems
	return &order, nil
}

func (s *OrderService) GetOrderByUser(userID string) ([]model.Order, error) {
	orders, err := s.Repo.OrderRepository.GetByUserID(userID)
	if err != nil {
		s.Logger.Error("error get order by id", zap.Error(err))
		return nil, err
	}

	var newOrders []model.Order
	for _, order := range orders {
		orderItems, err := s.GetOrderItems(order.ID)
		if err != nil {
			s.Logger.Error("error get order item by order id", zap.Error(err))
			return nil, err
		}
		order.OrderItems = orderItems
		newOrders = append(newOrders, order)
	}

	return newOrders, nil
}

func (s *OrderService) GetOrderItems(orderID int) ([]model.OrderItem, error) {
	orderItems, err := s.Repo.OrderRepository.GetOrderItems(orderID)
	if err != nil {
		s.Logger.Error("error get order item by order id", zap.Error(err))
		return nil, err
	}
	var newOrderItems []model.OrderItem
	for _, orderItem := range orderItems {
		product, err := s.Repo.ProductRepository.GetByID(orderItem.ProductID)
		if err != nil {
			s.Logger.Error("error get product by id", zap.Error(err))
			return nil, err
		}
		orderItem.Product = product
		NewVariants := []model.OrderItemVariant{}
		if product.HasVariant {
			for _, variantItem := range orderItem.Variants {
				variant, err := s.Repo.VariantRepository.GetVariantByID(int(variantItem.VariantID.Int64))
				if err != nil {
					s.Logger.Error("error get variant by id", zap.Error(err))
					return nil, err
				}
				option, err := s.Repo.VariantRepository.GetVariantOptionByID(int(variantItem.OptionID.Int64))
				if err != nil {
					s.Logger.Error("error get variant option by id", zap.Error(err))
					return nil, err
				}

				variantItem.Variant = variant
				variantItem.Option = option
				NewVariants = append(NewVariants, variantItem)
			}
			orderItem.Variants = NewVariants
		}
		newOrderItems = append(newOrderItems, orderItem)
	}
	return newOrderItems, nil
}
