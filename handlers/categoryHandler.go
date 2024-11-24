package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"

	// "github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewCategoryHandler(service service.MainService, log *zap.Logger) CategoryHandler {
	return CategoryHandler{Service: service, Logger: log}
}

func (h *CategoryHandler) GetAllCategoryHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "CategoryHandler"
	functionName := "GetAllCategoryHandler"

	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	var paginationInput model.Pagination
	page := r.URL.Query().Get("page")
	perPage := r.URL.Query().Get("perPage")

	// Parse pagination query parameters
	if page != "" {
		if parsedPage, err := strconv.Atoi(page); err == nil {
			paginationInput.Page = parsedPage
		} else {
			h.Logger.Warn("Invalid 'page' query parameter",
				zap.String("page", page),
				zap.String("handler", handlerName),
				zap.String("function", functionName),
			)
		}
	}

	if perPage != "" {
		if parsedPerPage, err := strconv.Atoi(perPage); err == nil {
			paginationInput.PerPage = parsedPerPage
		} else {
			h.Logger.Warn("Invalid 'perPage' query parameter",
				zap.String("perPage", perPage),
				zap.String("handler", handlerName),
				zap.String("function", functionName),
			)
		}
	}

	h.Logger.Info("Fetching categories",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.Int("page", paginationInput.Page),
		zap.Int("perPage", paginationInput.PerPage),
	)

	// Fetch categories and pagination details
	categories, pagination, err := h.Service.CategoryService.GetAllCategory(paginationInput)
	if err != nil {
		h.Logger.Error("Failed to fetch categories",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Int("page", paginationInput.Page),
			zap.Int("perPage", paginationInput.PerPage),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get all categories")
		return
	}

	// Calculate total pages
	TotalPage := 0
	if pagination.PerPage > 0 {
		TotalPage = (pagination.CountData + pagination.PerPage - 1) / pagination.PerPage // Round up
	}

	h.Logger.Info("Categories successfully retrieved",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.Int("page", pagination.Page),
		zap.Int("perPage", pagination.PerPage),
		zap.Int("countData", pagination.CountData),
		zap.Int("totalPage", TotalPage),
	)

	JsonResponse.SendPaginatedResponse(w, categories, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "Categories successfully retrieved")
}
