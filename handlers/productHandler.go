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
		h.Logger.Warn("Invalid HTTP method", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetAllProductHandler"))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, errMessage)
		return
	}

	var paginationInput model.Pagination
	var productFilter model.ProductDTO

	// Extract pagination and filters from query parameters
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

	// Fetch products
	h.Logger.Debug("Fetching products", zap.Any("filters", productFilter), zap.Any("pagination", paginationInput))
	products, pagination, err := h.Service.ProductService.GetAllProduct(productFilter, paginationInput)
	if err != nil {
		h.Logger.Error("Failed to fetch products", zap.Error(err), zap.String("handler", "Product"), zap.String("function", "GetAllProductHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get products")
		return
	}

	totalPage := 0
	if pagination.PerPage > 0 {
		totalPage = (pagination.CountData + pagination.PerPage - 1) / pagination.PerPage // Calculate total pages
	}

	h.Logger.Info("Products retrieved successfully", zap.Int("count", len(products)), zap.Int("page", pagination.Page), zap.Int("totalPages", totalPage))
	JsonResponse.SendPaginatedResponse(w, products, pagination.Page, pagination.PerPage, pagination.CountData, totalPage, "Products successfully retrieved")
}

func (h *ProductHandler) GetProductByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Warn("Invalid HTTP method", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, errMessage)
		return
	}

	// Extract product ID
	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Warn("Missing product ID in request URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Product ID is required")
		return
	}
	productID, err := strconv.Atoi(id)
	if err != nil || productID <= 0 {
		h.Logger.Warn("Invalid product ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	// Fetch product details
	h.Logger.Debug("Fetching product by ID", zap.Int("productID", productID))
	product, err := h.Service.ProductService.GetProductByID(productID)
	if err != nil {
		h.Logger.Error("Failed to fetch product by ID", zap.Error(err), zap.String("handler", "Product"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get product")
		return
	}

	h.Logger.Info("Product retrieved successfully", zap.Int("productID", productID))
	JsonResponse.SendSuccess(w, product, "Product successfully retrieved")
}

func (h *ProductHandler) GetWeeklyPromotionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Warn("Invalid HTTP method", zap.String("method", r.Method), zap.String("handler", "Product"), zap.String("function", "GetWeeklyPromotionsHandler"))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, errMessage)
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

	// Fetch weekly promotions
	h.Logger.Debug("Fetching weekly promotions", zap.Any("pagination", paginationInput))
	weeklyPromo, pagination, err := h.Service.ProductService.GetPromoWeekly(paginationInput)
	if err != nil {
		h.Logger.Error("Failed to fetch weekly promotions", zap.Error(err), zap.String("handler", "Product"), zap.String("function", "GetWeeklyPromotionsHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get weekly promotions")
		return
	}

	totalPage := 0
	if pagination.PerPage > 0 {
		totalPage = (pagination.CountData + pagination.PerPage - 1) / pagination.PerPage // Calculate total pages
	}

	h.Logger.Info("Weekly promotions retrieved successfully", zap.Int("count", len(weeklyPromo)), zap.Int("page", pagination.Page), zap.Int("totalPages", totalPage))
	JsonResponse.SendPaginatedResponse(w, weeklyPromo, pagination.Page, pagination.PerPage, pagination.CountData, totalPage, "Weekly Promo successfully retrieved")
}
