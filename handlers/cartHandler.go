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

	var cartInput model.CartDTO
	err := json.NewDecoder(r.Body).Decode(&cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.CartService.AddCart(user.ID, cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Cart"), zap.String("function", "AddCartHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add product to cart")
		return
	}

	JsonResponse.SendCreated(w, nil, "Product added to cart successfully")
}

func (h *CartHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {
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
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Cart"), zap.String("function", "GetCartHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get products in cart")
		return
	}
	JsonResponse.SendSuccess(w, cart, "Products in cart successfully retrieved")
}

func (h *CartHandler) UpdateCartHandler(w http.ResponseWriter, r *http.Request) {
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
		h.Logger.Error("Address ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	cartId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	var cartInput model.CartDTO
	err = json.NewDecoder(r.Body).Decode(&cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.CartService.UpdateCart(user.ID, cartId, cartInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Cart"), zap.String("function", "AddCartHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add product to cart")
		return
	}

	JsonResponse.SendSuccess(w, nil, "Product added to cart successfully")
}

func (h *CartHandler) DeleteProductInCartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only DELETE methods are allowed")
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
		h.Logger.Error("Product ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	cartID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"DeleteProductInCartHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.Service.CartService.DeleteProductInCart(cartID, user.ID)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"DeleteProductInCartHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to delete product from cart")
		return
	}

	JsonResponse.SendSuccess(w, nil, "Product deleted from cart successfully")
}

func (h *CartHandler) GetTotalCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	totalAmount, SubTotal, err := h.Service.CartService.GetTotalCart()
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Cart"), zap.String("function", "GetTotalCart"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get total cart amount")
		return
	}

	JsonResponse.SendSuccess(w, map[string]interface{}{
		"totalAmount": totalAmount,
		"subTotal":    SubTotal,
	}, "Total cart amount successfully retrieved")
}
