package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/middleware"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type CartHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewCartHandler(service service.MainService, log *zap.Logger) CartHandler {
	return CartHandler{Service: service, Logger: log}
}

func (h *CartHandler) AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only POST methods are allowed")
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

	var cartInput model.CartItemDTO
	err := json.NewDecoder(r.Body).Decode(&cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.CartService.AddProductToCart(user.ID, cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Cart"), zap.String("function", "AddCartHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add product to second cart")
		return
	}

	JsonResponse.SendCreated(w, nil, "Product added to second cart successfully")
}

func (h *CartHandler) UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only PUT methods are allowed")
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
	if id == "" {
		h.Logger.Error("Cart item ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Cart item ID is required")
		return
	}
	cartItemId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"UpdateCartItemHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}
	var cartItemInput model.CartItem
	err = json.NewDecoder(r.Body).Decode(&cartItemInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cartItemInput.ID = cartItemId

	err = h.Service.CartService.UpdateItemInCart(user.ID, cartItemInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"UpdateCartItemHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	JsonResponse.SendSuccess(w, nil, "Cart item updated successfully")
}

func (h *CartHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only DELETE methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	// Convert id to int
	cartItemId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "DeleteItem"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	err = h.Service.CartService.DeleteProductInCart(cartItemId)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "DeleteItem"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to delete cart item")
		return
	}
	JsonResponse.SendSuccess(w, nil, "Item deleted successfully")
}

func (h *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
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
	cart, err := h.Service.CartService.GetCartByUserID(user.ID)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "GetUserCart"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve user's cart")
		return
	}
	JsonResponse.SendSuccess(w, cart, "User's cart retrieved successfully")
}
