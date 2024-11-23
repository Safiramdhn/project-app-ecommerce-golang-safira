package repository

import (
	"database/sql"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type WishlistRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewWishlistRepository(db *sql.DB, logger *zap.Logger) WishlistRepository {
	return WishlistRepository{DB: db, Logger: logger}
}

func (repo *WishlistRepository) Create(wishlistInput model.WishlistDTO) error {
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

	sqlStatement := `INSERT INTO wishlist (user_id, product_id) VALUES ($1, $2)`
	repo.Logger.Info("Execute query", zap.String("query", sqlStatement), zap.String("Repository", "Wishlist"), zap.String("Function", "Create"))
	_, err = tx.Exec(sqlStatement, wishlistInput.UserID, wishlistInput.ProductID)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository",
			"Wishlist"), zap.String("Function", "Create"))
		return err
	}
	return nil
}

func (repo *WishlistRepository) GetAll(userID string, pagination model.Pagination) ([]model.Wishlist, model.Pagination, error) {
	var wishlist []model.Wishlist
	sqlStatement := `SELECT product_id FROM wishlist WHERE user_id = $1 AND status = 'active' LIMIT $2 OFFSET $3`
	limit := pagination.PerPage
	offset := (pagination.Page - 1) / limit

	rows, err := repo.DB.Query(sqlStatement, userID, limit, offset)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err), zap.String("Repository", "Wishlist"), zap.String("Function", "GetAll"))
		return nil, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int
		err = rows.Scan(&productID)
		if err != nil {
			repo.Logger.Error("Failed to scan row", zap.Error(err), zap.String("Repository", "Wishlist"), zap.String("Function", "GetAll"))
			return nil, pagination, err
		}
		wishlist = append(wishlist, model.Wishlist{UserID: userID, ProductID: productID})
	}

	totalCount, err := repo.CountWishlist(userID)
	if err != nil {
		return nil, pagination, err
	}
	pagination.CountData = totalCount
	return wishlist, pagination, nil
}

func (repo *WishlistRepository) Delete(userID string, wishlistID int) error {
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

	sqlStatement := `UPDATE wishlist SET status = 'deleted', deleted_at = NOW() WHERE id = $1 AND user_id = $2`

	repo.Logger.Info("Execute query", zap.String("query", sqlStatement), zap.String("Repository", "Wishlist"), zap.String("Function", "Delete"))
	_, err = tx.Exec(sqlStatement, wishlistID, userID)
	if err != nil {
		repo.Logger.Error("Failed to execute query", zap.Error(err))
		return err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Error("Failed to commit transaction", zap.Error(err), zap.String("Repository", "Wishlist"), zap.String("Function", "Delete"))
		return err
	}
	return nil
}

func (repo *WishlistRepository) CountWishlist(userID string) (int, error) {
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM wishlist WHERE user_id = $1 AND status = 'active';`

	repo.Logger.Info("Execute query", zap.String("query", countQuery), zap.String("Repository", "Wishlist"), zap.String("Function", "CountWishlist"))
	err := repo.DB.QueryRow(countQuery, userID).Scan(&totalCount)
	if err != nil {
		repo.Logger.Error("Error counting wishlist", zap.Error(err), zap.String("Repository", "Wishlist"), zap.String("Function", "CountWishlist"))
		return 0, err
	}

	return totalCount, nil
}
