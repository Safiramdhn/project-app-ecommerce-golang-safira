package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"go.uber.org/zap"
)

type CategoryHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	Config  util.Configuration
}

func NewCategoryHandler(service service.MainService, log *zap.Logger, config util.Configuration) CategoryHandler {
	return CategoryHandler{Service: service, Logger: log, Config: config}
}

func (h *CategoryHandler) GetAllCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Category"), zap.String("function", "GetAllCategoryHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	var paginationInput model.Pagination

	page := r.URL.Query().Get("page")
	if page != "" {
		paginationInput.Page, _ = strconv.Atoi(page)
	}
	perPage := r.URL.Query().Get("perPage")
	if perPage != "" {
		paginationInput.PerPage, _ = strconv.Atoi(perPage)
	}

	categories, pagination, err := h.Service.CategoryService.GetAllCategory(paginationInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Category"), zap.String("function", "GetAllCategoryHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get all categories")
		return
	}
	totalPage := pagination.CountData / pagination.PerPage
	JsonResponse.SendPaginatedResponse(w, categories, pagination.Page, pagination.PerPage, pagination.CountData, totalPage, "Categories successfully retrieved")
}
