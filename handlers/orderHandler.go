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

type OrderHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	// Config  util.Configuration
}

func NewOrderHandler(service service.MainService, log *zap.Logger) OrderHandler {
	return OrderHandler{Service: service, Logger: log}
}

func (h *OrderHandler) CreateOrderHanlder(w http.ResponseWriter, r *http.Request) {
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

	var orderInput model.OrderDTO
	err := json.NewDecoder(r.Body).Decode(&orderInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.OrderService.CreateOrder(user.ID, orderInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Order"), zap.String("function", "CreateOrderHanlder"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Order creation failed due internal error")
		return
	}
	JsonResponse.SendCreated(w, nil, "Order created successfully")
}

func (h *OrderHandler) GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.Service.OrderService.GetOrderByUser(user.ID)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Order"), zap.String("function", "GetOrderHistoryHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve order history")
		return
	}

	JsonResponse.SendSuccess(w, orders, "Order history successfully retrieved")
}

func (h *OrderHandler) GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Order ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Order"), zap.String("function", "GetOrderDetailsHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := h.Service.OrderService.GetOrderByID(orderID)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Order"), zap.String("function", "GetOrderDetailsHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve order details")
		return
	}
	JsonResponse.SendSuccess(w, order, "Order details successfully retrieved")
}
