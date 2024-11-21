package repository

import (
	"database/sql"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type VariantRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewVariantRepository(db *sql.DB, logger *zap.Logger) VariantRepository {
	return VariantRepository{DB: db, Logger: logger}
}

func (repo *VariantRepository) GetByProductId(productId int) ([]model.Variant, error) {
	sqlStatement := `SELECT id, attribute_name FROM variations WHERE product_id = $1 AND status = 'active'`
	repo.Logger.Info("Executing query", zap.String("query", sqlStatement), zap.String("Repository", "Variant"), zap.String("Function", "GetByProductId"))
	rows, err := repo.DB.Query(sqlStatement, productId)
	if err != nil {
		repo.Logger.Error("Error getting variants by product ID", zap.Error(err),
			zap.String("Repository", "Variant"),
			zap.String("Function", "GetByProductId"),
			zap.Duration("duration", time.Since(startTime)))

		return nil, err
	}
	defer rows.Close()

	var variants []model.Variant
	for rows.Next() {
		var variant model.Variant
		err := rows.Scan(&variant.ID, &variant.AttributeName)
		if err != nil {
			repo.Logger.Error("Error scanning variant row", zap.Error(err),
				zap.String("Repository", "Variant"),
				zap.String("Function", "GetByProductId"),
				zap.Duration("duration", time.Since(startTime)))
			return nil, err
		}
		variantOption, err := repo.GetVariantOptions(variant.ID)
		if err != nil {
			return nil, err
		}
		variant.VariantOption = append(variant.VariantOption, variantOption...)
		variants = append(variants, variant)
	}

	return variants, nil
}

func (repo *VariantRepository) GetVariantOptions(variantId int) ([]model.VariantOption, error) {
	sqlStatement := `SELECT id, option_value, additional_price, stock FROM variation_options WHERE variation_id = $1 AND status = 'active'`
	rows, err := repo.DB.Query(sqlStatement, variantId)
	if err != nil {
		repo.Logger.Error("Error getting variant options by variant ID", zap.Error(err),
			zap.String("Repository", "Variant"),
			zap.String("Function", "GetVariantOptions"),
			zap.Duration("duration", time.Since(startTime)))
		return nil, err
	}
	defer rows.Close()
	var variantOptions []model.VariantOption
	for rows.Next() {
		var variantOption model.VariantOption
		err := rows.Scan(&variantOption.ID, &variantOption.OptionValue, &variantOption.AdditionalPrice, &variantOption.Stock)
		if err != nil {
			repo.Logger.Error("Error scanning variant option row", zap.Error(err),
				zap.String("Repository", "Variant"),
				zap.String("Function", "GetVariantOption"),
				zap.Duration("duration", time.Since(startTime)))
			return nil, err
		}
		variantOptions = append(variantOptions, variantOption)
	}
	return variantOptions, nil
}
