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

type RecommendationHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewRecommendationHandler(service service.MainService, log *zap.Logger) RecommendationHandler {
	return RecommendationHandler{Service: service, Logger: log}
}

var recommedFilter = model.RecommendationDTO{
	IsRecommended: false,
	SetInBanner:   false,
}

func (h *RecommendationHandler) GetRecommendationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Recommendation"), zap.String("function", "GetRecommendations"))
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

	recommedFilter.IsRecommended = true
	recommendations, pagination, err := h.Service.RecommendationService.GetProductRecommendations(recommedFilter, paginationInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get product recommendations")
		return
	}
	if pagination.CountData/pagination.PerPage > 0 {
		TotalPage = pagination.CountData / pagination.PerPage
	}
	JsonResponse.SendPaginatedResponse(w, recommendations, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "Product recommendations successfully retrieved")
}

func (h *RecommendationHandler) GetBannerProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Recommendation"), zap.String("function", "GetBannerProduct"))
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

	recommedFilter.SetInBanner = true
	bannerProduct, pagination, err := h.Service.RecommendationService.GetProductRecommendations(recommedFilter, paginationInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get product recommendations")
		return
	}
	if pagination.CountData/pagination.PerPage > 0 {
		TotalPage = pagination.CountData / pagination.PerPage
	}
	JsonResponse.SendPaginatedResponse(w, bannerProduct, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "Banner products successfully retrieved")
}
