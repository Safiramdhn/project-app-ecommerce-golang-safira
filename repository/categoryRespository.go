package repository

import (
	"database/sql"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type CategoryRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewCategoryRepository(db *sql.DB, logger *zap.Logger) CategoryRepository {
	return CategoryRepository{DB: db, Logger: logger}
}

var startTime = time.Now()

func (repo CategoryRepository) GetAll(pagination model.Pagination) ([]model.Category, model.Pagination, error) {
	var categories []model.Category

	sqlStatement := `SELECT id, name FROM categories LIMIT $1 OFFSET $2`
	limit := pagination.PerPage
	offset := (pagination.Page - 1) / limit

	repo.Logger.Info("running query",
		zap.String("query", sqlStatement),
		zap.String("Repository", "Category"),
		zap.String("Function", "GetAll"),
	)
	rows, err := repo.DB.Query(sqlStatement, limit, offset)
	if err != nil {
		repo.Logger.Error("error when query database",
			zap.Error(err),
			zap.String("Repository", "Category"),
			zap.String("Function", "GetAll"),
			zap.Duration("duration", time.Since(startTime)))

		return categories, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			repo.Logger.Error("error when scanning row", zap.Error(err),
				zap.String("Repository", "Category"),
				zap.Duration("duration", time.Since(startTime)))

			return categories, pagination, err
		}
		categories = append(categories, category)
	}

	totalCount, err := repo.CountCategories()
	if err != nil {
		return categories, pagination, err
	}

	pagination.CountData = totalCount
	return categories, pagination, nil
}

func (repo CategoryRepository) CountCategories() (int, error) {
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM categories`

	repo.Logger.Info("running query", zap.String("query", countQuery), zap.String("Repository", "Category"), zap.String("Function", "CountCategories"))
	err := repo.DB.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		repo.Logger.Error("Error counting Category", zap.Error(err),
			zap.String("Repository", "Category"),
			zap.String("Function", "CountCategories"),
			zap.Duration("duration", time.Since(startTime)))
		return 0, err
	}

	return totalCount, nil
}
