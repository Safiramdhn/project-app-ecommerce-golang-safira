package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/middleware"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type WishlistHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewWishlistHandler(service service.MainService, log *zap.Logger) WishlistHandler {
	return WishlistHandler{Service: service, Logger: log}
}

var ContextKey = middleware.UserClaimsContextKey

func (h *WishlistHandler) AddWishlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error(errMessage, zap.String("method", r.Method), zap.String("handler", "Wishlist"), zap.String("function", "AddWishlistHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context")
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context")
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var wishlistInput model.WishlistDTO
	err := json.NewDecoder(r.Body).Decode(&wishlistInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	wishlistInput.UserID = user.ID

	err = h.Service.WishlistService.AddProductToWishlist(wishlistInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Wishlist"), zap.String("function", "AddWishlistHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add product to wishlist")
		return
	}

	JsonResponse.SendCreated(w, nil, "Product added to wishlist successfully")
}

func (h *WishlistHandler) GetWishlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context")
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context")
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	paginationInput := model.Pagination{}
	page := r.URL.Query().Get("page")
	if page != "" {
		paginationInput.Page, _ = strconv.Atoi(page)
	}
	perPage := r.URL.Query().Get("perPage")
	if perPage != "" {
		paginationInput.PerPage, _ = strconv.Atoi(perPage)
	}

	wishlist, pagination, err := h.Service.WishlistService.GetWishlistByUserId(user.ID, paginationInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Wishlist"), zap.String("function", "GetWishlistHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get wishlist")
		return
	}

	if pagination.CountData/pagination.PerPage > 0 {
		TotalPage = pagination.CountData / pagination.PerPage
	}
	JsonResponse.SendPaginatedResponse(w, wishlist, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "Wishlist successfully retrieved")
}

func (h *WishlistHandler) RemoveProductFromWishlistHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context")
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context")
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	id := chi.URLParam(r, "id")
	var wishlistID int
	if id != "" {
		wishlistID, _ = strconv.Atoi(id)
	}

	err := h.Service.WishlistService.RemoveProductFromWishlist(user.ID, wishlistID)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Wishlist"), zap.String("function", "RemoveProductFromWishlistHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to remove product from wishlist")
		return
	}
	JsonResponse.SendSuccess(w, nil, "Product removed from wishlist successfully")
}
