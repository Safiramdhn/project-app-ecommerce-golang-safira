package repository

import (
	"database/sql"

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

func (repo CategoryRepository) GetAll(pagination model.Pagination) ([]model.Category, model.Pagination, error) {
	var categories []model.Category

	sqlStatement := `SELECT id, name FROM categories LIMIT $1 OFFSET $2`
	limit := pagination.PerPage
	offset := (pagination.Page - 1) / limit

	rows, err := repo.DB.Query(sqlStatement, limit, offset)
	if err != nil {
		repo.Logger.Error("error when query database", zap.Error(err))
		return categories, pagination, err
	}
	defer rows.Close()

	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			repo.Logger.Error("error when scanning row", zap.Error(err))
			return categories, pagination, err
		}
		categories = append(categories, category)
	}

	err = repo.DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&pagination.CountData)
	if err != nil {
		return categories, pagination, err
	}

	return categories, pagination, nil
}
