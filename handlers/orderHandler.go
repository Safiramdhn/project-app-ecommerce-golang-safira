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

func (h *OrderHandler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "OrderHandler"
	functionName := "CreateOrderHandler"

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

	var orderInput model.OrderDTO
	err := json.NewDecoder(r.Body).Decode(&orderInput)
	if err != nil {
		h.Logger.Error("Failed to decode request body",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.OrderService.CreateOrder(user.ID, orderInput)
	if err != nil {
		h.Logger.Error("Failed to create order",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("userId", user.ID),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Order creation failed due to internal error")
		return
	}

	h.Logger.Info("Order created successfully",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.String("userId", user.ID),
	)
	JsonResponse.SendCreated(w, nil, "Order created successfully")
}

func (h *OrderHandler) GetOrderHistoryHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "OrderHandler"
	functionName := "GetOrderHistoryHandler"

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

	h.Logger.Info("Fetching order history",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.String("userId", user.ID),
	)

	orders, err := h.Service.OrderService.GetOrderByUser(user.ID)
	if err != nil {
		h.Logger.Error("Failed to retrieve order history",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("userId", user.ID),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve order history")
		return
	}

	h.Logger.Info("Order history retrieved successfully",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.String("userId", user.ID),
		zap.Int("orderCount", len(orders)),
	)
	JsonResponse.SendSuccess(w, orders, "Order history successfully retrieved")
}

func (h *OrderHandler) GetOrderDetailsHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "OrderHandler"
	functionName := "GetOrderDetailsHandler"

	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Order ID is missing in the URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Order ID is required")
		return
	}

	orderID, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid order ID format",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("orderID", id),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	h.Logger.Info("Fetching order details",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.Int("orderID", orderID),
	)

	order, err := h.Service.OrderService.GetOrderByID(orderID)
	if err != nil {
		h.Logger.Error("Failed to retrieve order details",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Int("orderID", orderID),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to retrieve order details")
		return
	}

	h.Logger.Info("Order details retrieved successfully",
		zap.String("handler", handlerName),
		zap.String("function", functionName),
		zap.Int("orderID", orderID),
	)
	JsonResponse.SendSuccess(w, order, "Order details successfully retrieved")
}
