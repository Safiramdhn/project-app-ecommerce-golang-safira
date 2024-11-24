package repository

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type OrderRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewOrderRepository(db *sql.DB, logger *zap.Logger) OrderRepository {
	return OrderRepository{DB: db, Logger: logger}
}

func (repo OrderRepository) Create(orderInput model.Order) (model.Order, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderInput, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO orders (user_id, address_id, total_amount, total_price, shipping_type, shipping_cost, payment_method) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = tx.QueryRow(sqlStatement, orderInput.UserID, orderInput.AddressID, orderInput.TotalAmount, orderInput.TotalPrice, orderInput.ShippingType, orderInput.ShippingCost, orderInput.PaymentMethod).Scan(&orderInput.ID)
	if err != nil {
		repo.Logger.Error("Failed to create order", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderInput, err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderInput, err
	}
	return orderInput, nil
}

func (repo OrderRepository) AddOrderItem(orderItemInput model.OrderItem) (model.OrderItem, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderItemInput, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO order_items (order_id, product_id, amount, subtotal) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(sqlStatement, orderItemInput.OrderID, orderItemInput.ProductID, orderItemInput.Amount, orderItemInput.SubTotal).Scan(&orderItemInput.ID)
	if err != nil {
		repo.Logger.Error("Failed to add order item", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderItemInput, err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return orderItemInput, err
	}
	return orderItemInput, nil
}

func (repo OrderRepository) AddOrderItemVariant(variantInput model.OrderItemVariant) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO order_item_variants (order_item_id, variant_id, option_id) VALUES ($1, $2, $3)`
	_, err = tx.Exec(sqlStatement, variantInput.OrderItemID, variantInput.VariantID, variantInput.OptionID)
	if err != nil {
		repo.Logger.Error("Failed to add order item variant", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return err
	}
	return nil
}

func (repo OrderRepository) UpdateOrderStatus(orderId int, orderStatus string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "Order"), zap.String("Function", "Create"))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `UPDATE orders SET order_status = $1 WHERE id = $2`
	_, err = tx.Exec(sqlStatement, orderStatus, orderId)
	if err != nil {
		repo.Logger.Error("Failed to update order status", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "Create"))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "Create"))
		return err
	}
	return nil
}

func (repo OrderRepository) GetByID(id int) (model.Order, error) {
	var order model.Order
	sqlStatement := `SELECT id, address_id, shipping_type, total_amount, total_price, order_status FROM orders WHERE id = $1`
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&order.ID, &order.AddressID, &order.ShippingType, &order.TotalAmount, &order.TotalPrice, &order.OrderStatus)
	if err == sql.ErrNoRows {
		return order, nil
	} else if err != nil {
		repo.Logger.Error("Failed to get order by ID", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "GetByID"))
		return order, err
	}

	return order, nil
}

func (repo OrderRepository) GetByUserID(userID string) ([]model.Order, error) {
	var order []model.Order
	sqlStatement := `SELECT id, total_amount, total_price, order_status FROM orders WHERE user_id = $1`
	rows, err := repo.DB.Query(sqlStatement, userID)
	if err != nil {
		repo.Logger.Error("Failed to get order by user ID", zap.Error(err), zap.String(
			"repository", "Order"), zap.String("Function", "GetByUserID"))
		return order, err
	}
	defer rows.Close()

	for rows.Next() {
		var o model.Order
		err := rows.Scan(&o.ID, &o.TotalAmount, &o.TotalPrice, &o.OrderStatus)
		if err != nil {
			repo.Logger.Error("Failed to scan order", zap.Error(err), zap.String("repository", "order"),
				zap.String("Function", "GetByUserID"))
			return nil, err
		}
		order = append(order, o)
	}

	return order, nil
}

func (repo OrderRepository) GetOrderItems(orderId int) ([]model.OrderItem, error) {
	var orderItems []model.OrderItem
	sqlStatement := `SELECT id, order_id, product_id, amount, subtotal FROM order_items WHERE order_id = $1`
	rows, err := repo.DB.Query(sqlStatement, orderId)
	if err != nil {
		repo.Logger.Error("Failed to get order items by order ID", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "GetOrderItems"))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderItem model.OrderItem
		err = rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Amount, &orderItem.SubTotal)
		if err != nil {
			repo.Logger.Error("Failed to scan order item", zap.Error(err), zap.String("repository", "order"),
				zap.String("Function", "GetOrderItems"))
			return nil, err
		}

		variant, err := repo.GetOrderItemVariants(orderItem.ID)
		if err != nil {
			return nil, err
		}

		orderItem.Variants = variant
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

func (repo OrderRepository) GetOrderItemVariants(itemId int) ([]model.OrderItemVariant, error) {
	var orderItemVariants []model.OrderItemVariant
	sqlStatement := `SELECT id, variant_id, option_id price FROM order_item_variants
	WHERE order_item_id = $1`
	rows, err := repo.DB.Query(sqlStatement, itemId)
	if err != nil {
		repo.Logger.Error("Failed to get order item variants by order item ID", zap.Error(err),
			zap.String("repository", "Order"), zap.String("Function", "GetOrderItemVariants"))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var orderItemVariant model.OrderItemVariant
		err = rows.Scan(&orderItemVariant.ID, &orderItemVariant.VariantID, &orderItemVariant.OptionID)
		if err != nil {
			repo.Logger.Error("Failed to scan order item variant", zap.Error(err),
				zap.String("repository", "Order"), zap.String("Function", "GetOrderItemVariants"))
			return nil, err
		}
	}
	return orderItemVariants, nil
}

func (repo OrderRepository) CountProduct(productID int) (int, error) {
	sqlStatement := `SELECT COUNT(*) FROM order_items WHERE product_id = $1`
	var count int
	err := repo.DB.QueryRow(sqlStatement, productID).Scan(&count)
	if err != nil {
		repo.Logger.Error("Failed to count products by ID", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "CountProduct"))
		return 0, err
	}
	return count, nil
}
