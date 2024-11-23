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

// Create a new cart
func (repo CartRepository) Create(userId string, cartInput model.Cart) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
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

	sqlStatement := `INSERT INTO cart (user_id, product_id, variant_ids, variant_option_ids, amount, total_price) VALUES ($1, $2, $3, $4)`

	repo.Logger.Info("Execute query", zap.String("query", sqlStatement), zap.String("Repository", "Cart"), zap.String("Function", "Create"))
	_, err = repo.DB.Exec(sqlStatement, userId, cartInput.ProductID, cartInput.VariantIDs, cartInput.VariantOptionIDs, cartInput.Amount, cartInput.TotalPrice)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository",
			"Cart"), zap.String("Function", "Create"))
		return err
	}
	return nil
}

func (repo CartRepository) Update(id int, userId string, cartInput model.Cart) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
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

	fields := map[string]interface{}{}

	if cartInput.Amount != 0 {
		fields["amount"] = cartInput.Amount
	}

	for _, v := range cartInput.VariantOptionIDs {
		if v.Valid && v.Int64 != 0 {
			fields["variant_option_ids"] = cartInput.VariantOptionIDs
		}
	}
	fields["total_price"] = cartInput.TotalPrice
	fields["updated_at"] = time.Now()

	// Build the SET clause dynamically
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

	queryStatement := `UPDATE cart SET ` + helper.JoinStrings(setClauses, ", ") + `
	WHERE id = $` + strconv.Itoa(index) +
		` AND user_id = $` + strconv.Itoa(index+1) +
		` AND status = 'active'`
	values = append(values, id, userId)

	repo.Logger.Info("Executing query", zap.Int("cart_id", id),
		zap.String("query", queryStatement), zap.String("repository", "Cart"),
		zap.String("function", "Update"))
	_, err = tx.Exec(queryStatement, values...)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("repository", "Cart"), zap.String("Function", "Update"))
		return err
	}
	return nil
}

func (repo CartRepository) Delete(id int, userId string) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Error("Failed to start transaction", zap.Error(err), zap.String("Repository", "User"), zap.String("Function", "Create"))
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-panic after rollback
		} else if err != nil {
			repo.Logger.Error("Error executing transaction", zap.Error(err), zap.String("repository", "User"), zap.String("Function", "Create"))
			tx.Rollback()
		}
	}()

	sqlStatement := `UPDATE cart SET status = 'deleted', deleted_at = NOW() WHERE id = $1 AND user_id = $2`
	repo.Logger.Info("Executing query", zap.Int("address_id", id),
		zap.String("query", sqlStatement), zap.String("repository", "Cart"),
		zap.String("function", "Delete"))
	_, err = tx.Exec(sqlStatement, id, userId)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (repo CartRepository) GetAll(userId string) ([]model.Cart, error) {
	sqlStatement := `SELECT id, product_id, variant_ids, variant_options, amount, total_price FROM cart WHERE user_id = $1 AND status = 'active'`
	repo.Logger.Info("Executing query", zap.String("query", sqlStatement), zap.String("repository", "Cart"), zap.String("Function", "GetAll"))

	rows, err := repo.DB.Query(sqlStatement, userId)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("repository", "Cart"))
		return nil, err
	}
	defer rows.Close()

	var carts []model.Cart
	for rows.Next() {
		var cart model.Cart
		err := rows.Scan(&cart.ID, &cart.ProductID, &cart.VariantIDs, &cart.VariantOptionIDs, &cart.Amount, &cart.TotalPrice)
		if err != nil {
			repo.Logger.Error("Failed to scan row", zap.Error(err), zap.String("repository", "Cart"), zap.String("Function", "GetAll"))
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func (repo CartRepository) CountItemCart() (int, float64, error) {
	sqlStatement := `SELECT SUM(amount), SUM(total_price) as total_price FROM cart WHERE status = 'active'`
	repo.Logger.Info("Executing query", zap.String("query", sqlStatement), zap.String("repository", "Cart"), zap.String("Function", "CountCart"))

	var totalItem int
	var totalPrice float64
	err := repo.DB.QueryRow(sqlStatement).Scan(&totalItem, &totalPrice)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("repository", "Cart"))
		return 0, 0, err
	}
	return totalItem, totalPrice, nil
}
