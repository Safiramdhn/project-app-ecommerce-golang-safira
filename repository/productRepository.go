package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
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
	sqlStatement := `SELECT id, name, description, price, discount, rating, photo_url, has_variant, total_stock FROM products WHERE id = $1 AND status = 'active'`

	repo.Logger.Info("running query", zap.String("query", sqlStatement), zap.String("Repository", "Product"), zap.String("Function", "GetByID"))
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Discount, &product.Rating, &product.PhotoURL, &product.HasVariant, &product.TotalStock)
	if err == sql.ErrNoRows {
		repo.Logger.Info("product not found",
			zap.Int("product id", id),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetByID"),
			zap.Duration("duration", time.Since(startTime)))

		return product, nil
	} else if err != nil {
		repo.Logger.Error("error getting product by id",
			zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetByID"),
			zap.Duration("duration", time.Since(startTime)))
		return product, err
	}

	product.PriceAfterDiscount = helper.CalculateDiscountPrice(product.Price, product.Discount)
	return product, nil
}

func (repo ProductRepository) GetAll(productFilter model.ProductDTO, pagination model.Pagination) ([]model.Product, model.Pagination, error) {
	var products []model.Product
	var filterArgs []interface{}

	// Build base SQL query
	sqlStatement := `
        SELECT id, name, description, price, discount, rating, photo_url, has_variant, total_stock 
        FROM products
        WHERE status = 'active'
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
		zap.String("Repository", "Product"),
		zap.String("function", "GetAllProducts"),
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
		repo.Logger.Error("Error retrieving products", zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetByID"),
			zap.Duration("duration", time.Since(startTime)))

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
			&product.TotalStock,
		); err != nil {
			repo.Logger.Error("Error scanning product", zap.Error(err),
				zap.String("Repository", "Product"),
				zap.String("Function", "GetAll"),
				zap.Duration("duration", time.Since(startTime)))

			return nil, pagination, err
		}
		isNewProduct, err := repo.GetNewProducts(product.ID)
		if err != nil {
			return nil, pagination, err
		}
		orderedProducts, err := repo.CountProductFromOrder(product.ID)
		if err != nil {
			return nil, pagination, err
		}
		product.PriceAfterDiscount = helper.CalculateDiscountPrice(product.Price, product.Discount)

		product.SpecialProduct.IsNewProduct = isNewProduct
		if orderedProducts > 10 {
			product.SpecialProduct.IsBestSelling = true
		} else {
			product.SpecialProduct.IsBestSelling = false
		}
		products = append(products, product)
	}

	// Check for errors during rows iteration
	if err := rows.Err(); err != nil {
		repo.Logger.Error("Error during rows iteration", zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetAll"),
			zap.Duration("duration", time.Since(startTime)))
		return nil, pagination, err
	}

	// Get total product count
	totalCount, err := repo.CountProducts(productFilter)
	if err != nil {
		return nil, pagination, err
	}
	pagination.CountData = totalCount

	return products, pagination, nil
}

func (repo ProductRepository) CountProducts(productFilter model.ProductDTO) (int, error) {
	// Base query
	countQuery := `SELECT COUNT(*) FROM products WHERE status = 'active'`
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
	repo.Logger.Info("Execute count query", zap.String("query", countQuery), zap.String("Repository", "Product"), zap.String("Function", "CountProducts"))
	err := repo.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		repo.Logger.Error("Error counting products", zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "CountProducts"),
			zap.Duration("duration", time.Since(startTime)))
		return 0, err
	}

	return totalCount, nil
}

func (repo ProductRepository) GetNewProducts(id int) (bool, error) {
	sqlStatement := `SELECT(created_at > NOW() - INTERVAL '30 days') AS is_new_product FROM products WHERE id = $1 AND status = 'active';`
	var isNewProduct bool

	repo.Logger.Info("run sql statement", zap.String("query", sqlStatement), zap.String("Repository", "Product"), zap.String("Function", "GetNewProduct"))
	err := repo.DB.QueryRow(sqlStatement, id).Scan(&isNewProduct)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		repo.Logger.Error("Error getting new product status", zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetNewProducts"),
			zap.Int("id", id),
			zap.Duration("duration", time.Since(startTime)))
		return false, err
	}

	return isNewProduct, nil
}

func (repo ProductRepository) GetWeeklyPromo(pagination model.Pagination) ([]model.WeeklyPromo, model.Pagination, error) {
	sqlStatement := `SELECT id, product_id, start_date, end_date, promo_discount 
			FROM weekly_promos WHERE status = 'active' 
			AND start_date <= CURRENT_DATE AND end_date >= date_trunc('week', CURRENT_DATE)`

	var weeklyPromos []model.WeeklyPromo

	repo.Logger.Info("run sql statement", zap.String("query", sqlStatement), zap.String("Repository", "Product"), zap.String("Function", "GetWeeklyPromo"))
	rows, err := repo.DB.Query(sqlStatement)
	if err != nil {
		repo.Logger.Error("Error getting weekly promo", zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetWeeklyPromo"),
			zap.Duration("duration", time.Since(startTime)))
		return nil, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var weeklyPromo model.WeeklyPromo
		err = rows.Scan(&weeklyPromo.ID, &weeklyPromo.ProductID, &weeklyPromo.StartDate, &weeklyPromo.EndDate, &weeklyPromo.PromoDiscount)
		if err != nil {
			repo.Logger.Error("Error scanning weekly promo", zap.Error(err),
				zap.String("Repository", "Product"),
				zap.String("Function", "GetWeeklyPromo"),
				zap.Duration("duration", time.Since(startTime)))
			return nil, pagination, err
		}
		weeklyPromos = append(weeklyPromos, weeklyPromo)
	}
	return weeklyPromos, pagination, nil
}

func (repo ProductRepository) GetPromoProduct(productId int) (model.WeeklyPromo, error) {
	var weeklyPromo model.WeeklyPromo
	sqlStatement := `SELECT id, start_date, end_date, promo_discount 
	FROM weekly_promos WHERE product_id = $1 AND status = 'active' 
	AND start_date <= CURRENT_DATE AND end_date >= date_trunc('week', CURRENT_DATE)`

	repo.Logger.Info("running query", zap.String("query", sqlStatement), zap.String("Repository", "Product"), zap.String("Function", "GetByID"))
	err := repo.DB.QueryRow(sqlStatement, productId).Scan(&weeklyPromo.ID, &weeklyPromo.StartDate, &weeklyPromo.EndDate, &weeklyPromo.PromoDiscount)
	if err == sql.ErrNoRows {
		repo.Logger.Info("product not found",
			zap.Int("product id", productId),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetByID"),
			zap.Duration("duration", time.Since(startTime)))

		return weeklyPromo, nil
	} else if err != nil {
		repo.Logger.Error("error getting product by id",
			zap.Error(err),
			zap.String("Repository", "Product"),
			zap.String("Function", "GetByID"),
			zap.Duration("duration", time.Since(startTime)))
		return weeklyPromo, err
	}
	return weeklyPromo, nil
}

func (repo ProductRepository) CountProductFromOrder(productID int) (int, error) {
	sqlStatement := `SELECT COUNT(*) FROM order_items WHERE product_id = $1`
	var count int
	err := repo.DB.QueryRow(sqlStatement, productID).Scan(&count)
	if err != nil {
		repo.Logger.Error("Failed to count products by ID", zap.Error(err), zap.String("repository", "Order"), zap.String("Function", "CountProduct"))
		return 0, err
	}
	return count, nil
}
