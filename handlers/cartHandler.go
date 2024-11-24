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
	handlerName := "CartHandler"
	functionName := "AddToCartHandler"

	if r.Method != http.MethodPost {
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only POST methods are allowed")
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	var cartInput model.CartItemDTO
	err := json.NewDecoder(r.Body).Decode(&cartInput)
	if err != nil {
		h.Logger.Error("Failed to decode request body",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.CartService.AddProductToCart(user.ID, cartInput)
	if err != nil {
		h.Logger.Error("Failed to add product to cart",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add product to cart")
		return
	}

	h.Logger.Info("Product added to cart successfully",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendCreated(w, nil, "Product added to cart successfully")
}

func (h *CartHandler) UpdateCartItemHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "CartHandler"
	functionName := "UpdateCartItemHandler"

	if r.Method != http.MethodPut {
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only PUT methods are allowed")
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Cart item ID is missing in URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Cart item ID is required")
		return
	}

	cartItemId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse cart item ID",
			zap.String("id", id),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	var cartItemInput model.CartItem
	err = json.NewDecoder(r.Body).Decode(&cartItemInput)
	if err != nil {
		h.Logger.Error("Failed to decode request body",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	cartItemInput.ID = cartItemId

	err = h.Service.CartService.UpdateItemInCart(user.ID, cartItemInput)
	if err != nil {
		h.Logger.Error("Failed to update cart item",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to update cart item")
		return
	}

	h.Logger.Info("Cart item updated successfully",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendSuccess(w, nil, "Cart item updated successfully")
}

func (h *CartHandler) DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "CartHandler"
	functionName := "DeleteItemHandler"

	if r.Method != http.MethodDelete {
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only DELETE methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Cart item ID is missing in URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Cart item ID is required")
		return
	}

	cartItemId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse cart item ID",
			zap.String("id", id),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	err = h.Service.CartService.DeleteProductInCart(cartItemId)
	if err != nil {
		h.Logger.Error("Failed to delete cart item",
			zap.Int("cartItemId", cartItemId),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to delete cart item")
		return
	}

	h.Logger.Info("Cart item deleted successfully",
		zap.Int("cartItemId", cartItemId),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendSuccess(w, nil, "Item deleted successfully")
}

func (h *CartHandler) GetUserCart(w http.ResponseWriter, r *http.Request) {
	handlerName := "CartHandler"
	functionName := "GetUserCart"

	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == nil {
		h.Logger.Error("User ID not found in context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	h.Logger.Info("Fetching user cart",
		zap.String("userID", user.ID),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)

	cart, err := h.Service.CartService.GetCartByUserID(user.ID)
	if err != nil {
		h.Logger.Error("Failed to retrieve user's cart",
			zap.String("userID", user.ID),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve user's cart")
		return
	}

	h.Logger.Info("User's cart retrieved successfully",
		zap.String("userID", user.ID),
		zap.Int("cartItemsCount", len(cart.Items)),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendSuccess(w, cart, "User's cart retrieved successfully")
}
