package repository

import (
	"database/sql"
	"fmt"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type ProductRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewProductRepository(db *sql.DB, logger *zap.Logger) ProductRepository {
	return ProductRepository{DB: db, Logger: logger}
}

func (repo ProductRepository) GetByID(id int) (model.Product, error) {
	var product model.Product
	sqlStatement := `SELECT id, name, description, price, discount, rating, photo_url, has_variant FROM products WHERE id = $1`
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Discount, &product.Rating, &product.PhotoURL, &product.HasVariant)
	if err == sql.ErrNoRows {
		repo.Logger.Info("product not found", zap.Int("product id", id))
		return product, nil
	} else if err != nil {
		repo.Logger.Error("error getting product by id", zap.Error(err))
		return product, err
	}

	return product, nil
}

func (repo ProductRepository) GetAll(productFilter model.ProductDTO, pagination model.Pagination) ([]model.Product, model.Pagination, error) {
	var products []model.Product
	var filterArgs []interface{}

	// Build base SQL query
	sqlStatement := `
        SELECT id, name, description, price, discount, rating, photo_url, has_variant 
        FROM products
        WHERE 1=1
    `

	// Add filters if provided
	if productFilter.Name != "" {
		sqlStatement += " AND name LIKE $1"
		filterArgs = append(filterArgs, "%"+productFilter.Name+"%")
	}

	if productFilter.CategoryID != 0 {
		sqlStatement += " AND category_id = $" + fmt.Sprint(len(filterArgs)+1)
		filterArgs = append(filterArgs, productFilter.CategoryID)
	}

	// Add pagination
	sqlStatement += " LIMIT $" + fmt.Sprint(len(filterArgs)+1) + " OFFSET $" + fmt.Sprint(len(filterArgs)+2)
	filterArgs = append(filterArgs, pagination.PerPage, (pagination.Page-1)*pagination.PerPage)

	// Log SQL statement and parameters (be cautious with sensitive data in production)
	repo.Logger.Info("Run Get All Products",
		zap.String("Repository", "ProductRepository"),
		zap.String("statement", sqlStatement),
		zap.Int("args_count", len(filterArgs)),
		// Optionally, mask sensitive data in args for logging purposes
		zap.Any("args", filterArgs),
	)

	// Execute query
	rows, err := repo.DB.Query(sqlStatement, filterArgs...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, pagination, nil
		}
		repo.Logger.Error("Error retrieving products", zap.Error(err))
		return nil, pagination, err
	}
	defer rows.Close()

	// Iterate and scan products
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Discount,
			&product.Rating,
			&product.PhotoURL,
			&product.HasVariant,
		); err != nil {
			repo.Logger.Error("Error scanning product", zap.Error(err))
			return nil, pagination, err
		}
		products = append(products, product)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		repo.Logger.Error("Error during rows iteration", zap.Error(err))
		return nil, pagination, err
	}

	// Get total product count
	totalCount, err := repo.CountProducts(productFilter)
	if err != nil {
		repo.Logger.Error("Error retrieving product count", zap.Error(err))
		return nil, pagination, err
	}
	pagination.CountData = totalCount

	return products, pagination, nil
}

func (repo ProductRepository) CountProducts(productFilter model.ProductDTO) (int, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`
	countArgs := []interface{}{}
	countArgIndex := 1

	// Add filters if provided
	if productFilter.Name != "" {
		countQuery += ` AND name ILIKE $` + fmt.Sprint(countArgIndex)
		countArgs = append(countArgs, "%"+productFilter.Name+"%")
		countArgIndex++
	}

	if productFilter.CategoryID != 0 {
		countQuery += ` AND category_id = $` + fmt.Sprint(countArgIndex)
		countArgs = append(countArgs, productFilter.CategoryID)
		countArgIndex++
	}

	// Execute count query
	var totalCount int
	err := repo.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		repo.Logger.Error("Error counting products", zap.Error(err))
		return 0, err
	}

	return totalCount, nil
}
