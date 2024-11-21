package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"go.uber.org/zap"
)

type RecommendationRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func NewRecommendationRepository(db *sql.DB, logger *zap.Logger) RecommendationRepository {
	return RecommendationRepository{DB: db, Logger: logger}
}

func (repo *RecommendationRepository) GetRecommendations(recommendFilter model.RecommendationDTO, pagination model.Pagination) ([]model.Recommendation, model.Pagination, error) {
	var recommendations []model.Recommendation
	var filterArgs []interface{}
	var argIndex = 1

	sqlStatement := `SELECT r.id, p.id, p.name, r.photo_url, r.is_recommended, r.set_in_banner,
				r.title, r.subtitle FROM recommendations r 
				JOIN products p ON p.id = r.product_id 
				WHERE p.status = 'active'`

	if recommendFilter.IsRecommended {
		sqlStatement += ` AND is_recommended = $` + fmt.Sprint(argIndex)
		filterArgs = append(filterArgs, recommendFilter.IsRecommended)
		argIndex++
	}
	if recommendFilter.SetInBanner {
		sqlStatement += ` AND set_in_banner = $` + fmt.Sprint(argIndex)
		filterArgs = append(filterArgs, recommendFilter.SetInBanner)
		argIndex++
	}

	sqlStatement += " LIMIT $" + fmt.Sprint(len(filterArgs)+1) + " OFFSET $" + fmt.Sprint(len(filterArgs)+2)
	filterArgs = append(filterArgs, pagination.PerPage, (pagination.Page-1)*pagination.PerPage)

	rows, err := repo.DB.Query(sqlStatement, filterArgs...)
	if err != nil {
		repo.Logger.Error("error getting recommendations", zap.Error(err))
		return nil, pagination, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		repo.Logger.Error("Error during rows iteration", zap.Error(err))
		return nil, pagination, err
	}

	for rows.Next() {
		var product model.Product
		var recommendation model.Recommendation

		if err := rows.Scan(&recommendation.ID,
			&product.ID,
			&product.Name,
			&recommendation.PhotoUrl,
			&recommendation.IsRecommended,
			&recommendation.SetInBanner,
			&recommendation.Title,
			&recommendation.Subtitle); err != nil {
			repo.Logger.Error("error scanning rows", zap.Error(err))
			return nil, pagination, err
		}

		recommendation.Product = product
		recommendation.PathUrl = fmt.Sprintf("/api/products/%d", product.ID)
		recommendations = append(recommendations, recommendation)
	}

	totalCount, err := repo.CountRecommendations(recommendFilter)
	if err != nil {
		return nil, pagination, err
	}
	pagination.CountData = totalCount
	return recommendations, pagination, nil
}

func (repo RecommendationRepository) CountRecommendations(recommendFilter model.RecommendationDTO) (int, error) {
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM recommendations WHERE status = 'active'`
	countArgs := []interface{}{}
	countArgIndex := 1

	if recommendFilter.IsRecommended {
		countQuery += ` AND is_recommended =
		 $` + fmt.Sprint(countArgIndex)
		countArgs = append(countArgs, recommendFilter.IsRecommended)
		countArgIndex++
	}
	if recommendFilter.SetInBanner {
		countQuery += ` AND set_in_banner = $` + fmt.Sprint(countArgIndex)
		countArgs = append(countArgs, recommendFilter.SetInBanner)
		countArgIndex++
	}

	repo.Logger.Info("running query", zap.String("query", countQuery),
		zap.String("Repository", "Recommendation"),
		zap.String("Function", "CountCategories"),
		zap.Int("args_count", len(countArgs)),
		// Optionally, mask sensitive data in args for logging purposes
		zap.Any("args", countArgs))

	err := repo.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		repo.Logger.Error("Error counting recommendation", zap.Error(err),
			zap.String("Repository", "Recommendation"),
			zap.String("Function", "CountRecommendations"),
			zap.Duration("duration", time.Since(startTime)))
		return 0, err
	}

	return totalCount, nil
}
