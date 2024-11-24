package repository

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type CartRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewCartRepository(db *sql.DB, logger *zap.Logger) CartRepository {
	return CartRepository{DB: db, Logger: logger}
}

func (repo CartRepository) Create(userID string) (model.Cart, error) {
	var Cart model.Cart
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Create"),
			zap.String("userID", userID),
		)
		return Cart, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic during transaction",
				zap.String("repository", "Cart"),
				zap.String("function", "Create"),
				zap.String("userID", userID),
				zap.Any("panic", p),
			)
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback()
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("repository", "Cart"),
				zap.String("function", "Create"),
				zap.String("userID", userID),
			)
		} else {
			repo.Logger.Info("Transaction completed successfully",
				zap.String("repository", "Cart"),
				zap.String("function", "Create"),
				zap.String("userID", userID),
			)
		}
	}()

	sqlStatement := `INSERT INTO carts (user_id) VALUES ($1) RETURNING id`
	err = tx.QueryRow(sqlStatement, userID).Scan(&Cart.ID)
	if err != nil {
		repo.Logger.Error("Failed to create cart",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Create"),
			zap.String("userID", userID),
		)
		return Cart, err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Create"),
			zap.String("userID", userID),
		)
		return Cart, err
	}
	return Cart, nil
}

func (repo CartRepository) Update(cartInput model.Cart) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Update"),
			zap.Int("cartID", cartInput.ID),
		)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic during transaction",
				zap.String("repository", "Cart"),
				zap.String("function", "Update"),
				zap.Int("cartID", cartInput.ID),
				zap.Any("panic", p),
			)
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback()
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("repository", "Cart"),
				zap.String("function", "Update"),
				zap.Int("cartID", cartInput.ID),
			)
		} else {
			repo.Logger.Info("Transaction completed successfully",
				zap.String("repository", "Cart"),
				zap.String("function", "Update"),
				zap.Int("cartID", cartInput.ID),
			)
		}
	}()

	sqlStatement := `UPDATE carts SET total_amount = $1, total_price = $2, updated_at = NOW() WHERE id = $3`
	_, err = tx.Exec(sqlStatement, cartInput.TotalAmount, cartInput.TotalPrice, cartInput.ID)
	if err != nil {
		repo.Logger.Error("Failed to update cart",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Update"),
			zap.Int("cartID", cartInput.ID),
		)
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "Update"),
			zap.Int("cartID", cartInput.ID),
		)
		return err
	}
	return nil
}

func (repo CartRepository) AddItem(itemInput model.CartItem) (model.CartItem, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "AddItem"),
			zap.Int("cartID", itemInput.CartID),
			zap.Int("productID", itemInput.ProductID),
		)
		return itemInput, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic during transaction",
				zap.String("repository", "Cart"),
				zap.String("function", "AddItem"),
				zap.Int("cartID", itemInput.CartID),
				zap.Int("productID", itemInput.ProductID),
				zap.Any("panic", p),
			)
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback()
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("repository", "Cart"),
				zap.String("function", "AddItem"),
				zap.Int("cartID", itemInput.CartID),
				zap.Int("productID", itemInput.ProductID),
			)
		} else {
			repo.Logger.Info("Transaction completed successfully",
				zap.String("repository", "Cart"),
				zap.String("function", "AddItem"),
				zap.Int("cartID", itemInput.CartID),
				zap.Int("productID", itemInput.ProductID),
			)
		}
	}()

	sqlStatement := `INSERT INTO cart_items (cart_id, product_id, amount, sub_total) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(sqlStatement, itemInput.CartID, itemInput.ProductID, itemInput.Amount, itemInput.SubTotal).Scan(&itemInput.ID)
	if err != nil {
		repo.Logger.Error("Failed to add cart item",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "AddItem"),
			zap.Int("cartID", itemInput.CartID),
			zap.Int("productID", itemInput.ProductID),
		)
		return itemInput, err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("repository", "Cart"),
			zap.String("function", "AddItem"),
			zap.Int("cartID", itemInput.CartID),
			zap.Int("productID", itemInput.ProductID),
		)
		return itemInput, err
	}
	return itemInput, nil
}

func (repo CartRepository) DeleteItem(itemID int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic",
				zap.String("Repository", "Cart"),
				zap.String("Function", "DeleteItem"),
				zap.Int("ItemID", itemID),
				zap.Any("Panic", p))
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("Repository", "Cart"),
				zap.String("Function", "DeleteItem"),
				zap.Int("ItemID", itemID))
			tx.Rollback()
		}
	}()

	var item model.CartItem
	sqlStatement := `SELECT cart_id, amount, sub_total FROM cart_items WHERE id = $1`
	err = tx.QueryRow(sqlStatement, itemID).Scan(&item.CartID, &item.Amount, &item.SubTotal)
	if err == sql.ErrNoRows {
		repo.Logger.Warn("Item not found",
			zap.Int("ItemID", itemID),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"))
		return nil
	}
	if err != nil {
		repo.Logger.Error("Failed to retrieve item details",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}

	cart, err := repo.GetByID(item.CartID)
	if err != nil {
		repo.Logger.Error("Failed to retrieve cart details",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}

	cart.TotalAmount -= item.Amount
	cart.TotalPrice -= item.SubTotal
	err = repo.Update(cart)
	if err != nil {
		repo.Logger.Error("Failed to update cart",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}

	sqlStatement = `UPDATE cart_items SET status = 'deleted', deleted_at = NOW() WHERE id = $1`
	_, err = tx.Exec(sqlStatement, itemID)
	if err != nil {
		repo.Logger.Error("Failed to delete cart item",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "DeleteItem"),
			zap.Int("ItemID", itemID))
		return err
	}
	return nil
}

func (repo CartRepository) AddItemVariant(cartItemID int, variantInput model.CartItemVariantDTO) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "AddItemVariant"),
			zap.Int("CartItemID", cartItemID))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic",
				zap.String("Repository", "Cart"),
				zap.String("Function", "AddItemVariant"),
				zap.Int("CartItemID", cartItemID),
				zap.Any("Panic", p))
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("Repository", "Cart"),
				zap.String("Function", "AddItemVariant"),
				zap.Int("CartItemID", cartItemID))
			tx.Rollback()
		}
	}()

	sqlStatement := `INSERT INTO cart_item_variants (cart_item_id, item_variant_id, option_id, additional_price) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(sqlStatement, cartItemID, variantInput.VariantID, variantInput.VariantOptionID, variantInput.AdditionalPrice)
	if err != nil {
		repo.Logger.Error("Failed to add cart item variant",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "AddItemVariant"),
			zap.Int("CartItemID", cartItemID))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "AddItemVariant"),
			zap.Int("CartItemID", cartItemID))
		return err
	}
	return nil
}

func (repo CartRepository) UpdateItem(itemInput model.CartItem) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "UpdateItem"),
			zap.Int("CartItemID", itemInput.ID))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Recovered from panic",
				zap.String("Repository", "Cart"),
				zap.String("Function", "UpdateItem"),
				zap.Int("CartItemID", itemInput.ID),
				zap.Any("Panic", p))
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction",
				zap.Error(err),
				zap.String("Repository", "Cart"),
				zap.String("Function", "UpdateItem"),
				zap.Int("CartItemID", itemInput.ID))
			tx.Rollback()
		}
	}()

	fields := map[string]interface{}{}
	if itemInput.Amount != 0 {
		fields["amount"] = itemInput.Amount
	}
	if itemInput.SubTotal != 0 {
		fields["sub_total"] = itemInput.SubTotal
	}
	fields["updated_at"] = time.Now()

	setClauses := []string{}
	values := []interface{}{}
	index := 1
	for field, value := range fields {
		setClauses = append(setClauses, field+"=$"+strconv.Itoa(index))
		values = append(values, value)
		index++
	}

	// Check if there are fields to update
	if len(setClauses) == 0 {
		return errors.New("no fields to update")
	}

	// Build the final query
	queryStatement := `
		UPDATE cart_items
		SET ` + helper.JoinStrings(setClauses, ", ") + `
		WHERE id = $` + strconv.Itoa(index) +
		` AND status = 'active'`
	values = append(values, itemInput.ID)

	repo.Logger.Info("Executing query",
		zap.Int("CartItemID", itemInput.ID),
		zap.String("Query", queryStatement),
		zap.String("Repository", "Cart"),
		zap.String("Function", "UpdateItem"))

	// Execute the query
	_, err = tx.Exec(queryStatement, values...)
	if err != nil {
		repo.Logger.Error("Failed to update cart item",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "UpdateItem"),
			zap.Int("CartItemID", itemInput.ID))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction",
			zap.Error(err),
			zap.String("Repository", "Cart"),
			zap.String("Function", "UpdateItem"),
			zap.Int("CartItemID", itemInput.ID))
		return err
	}
	return nil
}

func (repo CartRepository) GetByUserID(userID string) (model.Cart, error) {
	var result model.Cart
	sqlStatement := `SELECT id, total_amount, total_price FROM carts WHERE user_id = $1 AND status = 'active' AND cart_status = 'active'`
	err := repo.DB.QueryRow(sqlStatement, userID).Scan(&result.ID, &result.TotalAmount, &result.TotalPrice)
	if err == sql.ErrNoRows {
		repo.Logger.Info("No active cart found for user", zap.String("userID", userID), zap.String("repository", "Cart"), zap.String("function", "GetByUserID"))
		return result, nil
	} else if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("userID", userID), zap.String("repository", "Cart"), zap.String("function", "GetByUserID"))
		return result, err
	}
	repo.Logger.Info("Successfully retrieved cart", zap.String("userID", userID), zap.String("repository", "Cart"), zap.String("function", "GetByUserID"))
	return result, nil
}

func (repo CartRepository) GetByID(id int) (model.Cart, error) {
	var result model.Cart
	sqlStatement := `SELECT id, user_id, total_amount, total_price FROM carts WHERE id = $1 AND status = 'active' AND cart_status = 'active'`
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&result.ID, &result.UserID, &result.TotalAmount, &result.TotalPrice)
	if err == sql.ErrNoRows {
		repo.Logger.Info("No active cart found for ID", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetByID"))
		return result, nil
	} else if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetByID"))
		return result, err
	}
	repo.Logger.Info("Successfully retrieved cart by ID", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetByID"))
	return result, nil
}

func (repo CartRepository) GetItems(cartId int) ([]model.CartItem, error) {
	var cartItems []model.CartItem
	sqlStatement := `SELECT id, product_id, amount, sub_total FROM cart_items WHERE cart_id = $1`
	rows, err := repo.DB.Query(sqlStatement, cartId)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("cartId", cartId), zap.String("repository", "Cart"), zap.String("function", "GetItems"))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.CartItem
		err = rows.Scan(&item.ID, &item.ProductID, &item.Amount, &item.SubTotal)
		if err != nil {
			repo.Logger.Error("Failed to scan row", zap.Error(err), zap.Int("cartId", cartId), zap.String("repository", "Cart"), zap.String("function", "GetItems"))
			return nil, err
		}

		variant, err := repo.GetItemVariants(item.ID)
		if err != nil {
			repo.Logger.Error("Failed to get item variants", zap.Error(err), zap.Int("cartItemId", item.ID), zap.String("repository", "Cart"), zap.String("function", "GetItems"))
			return nil, err
		}
		item.ItemVariant = variant
		cartItems = append(cartItems, item)
	}
	repo.Logger.Info("Successfully retrieved cart items", zap.Int("cartId", cartId), zap.String("repository", "Cart"), zap.String("function", "GetItems"))
	return cartItems, nil
}

func (repo CartRepository) RecalculateTotal(cartID int) error {
	sqlStatement := `SELECT SUM(amount) as total_amount, SUM(sub_total) as total_price FROM cart_items WHERE cart_id = $1 AND status ='active'`
	var totalAmount, totalPrice float64
	err := repo.DB.QueryRow(sqlStatement, cartID).Scan(&totalAmount, &totalPrice)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("cartID", cartID), zap.String("repository", "Cart"), zap.String("function", "RecalculateTotal"))
		return err
	}

	cartInput := model.Cart{
		ID:          cartID,
		TotalAmount: int(totalAmount),
		TotalPrice:  totalPrice,
	}
	err = repo.Update(cartInput)
	if err != nil {
		repo.Logger.Error("Failed to update cart", zap.Error(err), zap.Int("cartID", cartID), zap.String("repository", "Cart"), zap.String("function", "RecalculateTotal"))
		return err
	}
	repo.Logger.Info("Successfully recalculated total", zap.Int("cartID", cartID), zap.String("repository", "Cart"), zap.String("function", "RecalculateTotal"))
	return nil
}

func (repo CartRepository) GetItemByID(id int) (model.CartItem, error) {
	var result model.CartItem
	sqlStatement := `SELECT id, cart_id, product_id, amount, sub_total FROM cart_items WHERE id = $1`
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&result.ID, &result.CartID, &result.ProductID, &result.Amount, &result.SubTotal)
	if err == sql.ErrNoRows {
		repo.Logger.Info("No cart item found by ID", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetItemByID"))
		return result, nil
	} else if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetItemByID"))
		return result, err
	}
	repo.Logger.Info("Successfully retrieved cart item", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "GetItemByID"))
	return result, nil
}

func (repo CartRepository) GetItemVariants(itemID int) ([]model.CarttemVariant, error) {
	var result []model.CarttemVariant
	sqlStatement := `SELECT id, cart_item_id, item_variant_id, option_id, additional_price  FROM cart_item_variants WHERE cart_item_id = $1`
	rows, err := repo.DB.Query(sqlStatement, itemID)
	if err == sql.ErrNoRows {
		repo.Logger.Info("No item variants found for cart item", zap.Int("itemID", itemID), zap.String("repository", "Cart"), zap.String("function", "GetItemVariants"))
		return result, nil
	} else if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("itemID", itemID), zap.String("repository", "Cart"), zap.String("function", "GetItemVariants"))
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var itemVariant model.CarttemVariant
		err = rows.Scan(&itemVariant.ID, &itemVariant.CartItemID, &itemVariant.VariantID, &itemVariant.OptionID, &itemVariant.AdditionalPrice)
		if err != nil {
			repo.Logger.Error("Failed to scan row", zap.Error(err), zap.Int("itemID", itemID), zap.String("repository", "Cart"), zap.String("function", "GetItemVariants"))
			return result, err
		}
		result = append(result, itemVariant)
	}
	repo.Logger.Info("Successfully retrieved item variants", zap.Int("itemID", itemID), zap.String("repository", "Cart"), zap.String("function", "GetItemVariants"))
	return result, nil
}

func (repo CartRepository) UpdateCartStatus(id int) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
		return err
	}
	repo.Logger.Info("Started transaction to update cart status", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			repo.Logger.Error("Panic occurred, rolling back transaction", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
			panic(p) // Re-panic after rollback
		} else if err != nil {
			tx.Rollback()
			repo.Logger.Error("Error executing transaction, rolling back", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
		}
	}()

	sqlStatement := `UPDATE carts SET cart_status = 'checkout' WHERE id = $1`
	_, err = tx.Exec(sqlStatement, id)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
		return err
	}
	repo.Logger.Info("Successfully executed cart status update query", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))

	if err = tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
		return err
	}
	repo.Logger.Info("Successfully committed transaction, cart status updated", zap.Int("id", id), zap.String("repository", "Cart"), zap.String("function", "UpdateCartStatus"))
	return nil
}
