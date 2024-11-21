package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	// "github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type ProductHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewProductHandler(service service.MainService, log *zap.Logger) ProductHandler {
	return ProductHandler{Service: service, Logger: log}
}

var TotalPage = 1

func (h *ProductHandler) GetAllProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetAllProductHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	var paginationInput model.Pagination
	var productFilter model.ProductDTO

	page := r.URL.Query().Get("page")
	if page != "" {
		paginationInput.Page, _ = strconv.Atoi(page)
	}
	perPage := r.URL.Query().Get("perPage")
	if perPage != "" {
		paginationInput.PerPage, _ = strconv.Atoi(perPage)
	}

	productName := r.URL.Query().Get("name")
	if productName != "" {
		productFilter.Name = productName
	}

	categoryID := r.URL.Query().Get("categoryId")
	if categoryID != "" {
		productFilter.CategoryID, _ = strconv.Atoi(categoryID)

	}

	products, pagination, err := h.Service.ProductService.GetAllProduct(productFilter, paginationInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetAllProductHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get products")
		return
	}
	if pagination.CountData/pagination.PerPage > 0 {
		TotalPage = pagination.CountData / pagination.PerPage
	}
	JsonResponse.SendPaginatedResponse(w, products, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "Products successfully retrieved")
}

func (h *ProductHandler) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	var productId int
	id := chi.URLParam(r, "id")
	if id != "" {
		productId, _ := strconv.Atoi(id)
		if productId <= 0 {
			errMessage := fmt.Sprintf("Invalid product id %s", id)
			h.Logger.Error("Invalid requested product ID", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetProductByIdHandler"))
			JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
			return
		}
	}

	product, err := h.Service.ProductService.GetProductByID(productId)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}

	JsonResponse.SendSuccess(w, product, "Product successfully retrieved")
}
