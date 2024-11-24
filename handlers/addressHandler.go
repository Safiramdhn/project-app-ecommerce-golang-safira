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

type AddressHandler struct {
	Service service.MainService
	Logger  *zap.Logger
}

func NewAddressHandler(service service.MainService, log *zap.Logger) AddressHandler {
	return AddressHandler{Service: service, Logger: log}
}

func (h *AddressHandler) AddAddressHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "AddressHandler"
	functionName := "AddAddressHandler"

	if r.Method != http.MethodPost {
		h.Logger.Error("Invalid method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only POST methods are allowed")
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

	var addressInput model.AddressDTO
	if err := json.NewDecoder(r.Body).Decode(&addressInput); err != nil {
		h.Logger.Error("Failed to decode JSON body",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Service.AddressService.AddAddress(user.ID, addressInput); err != nil {
		h.Logger.Error("Failed to add address",
			zap.Error(err),
			zap.String("userId", user.ID),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to add address")
		return
	}

	h.Logger.Info("Address successfully created",
		zap.String("userId", user.ID),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendCreated(w, nil, "Address successfully created")
}

func (h *AddressHandler) UpdateAddressHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "AddressHandler"
	functionName := "UpdateAddressHandler"

	if r.Method != http.MethodPut {
		h.Logger.Error("Invalid method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only PUT methods are allowed")
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

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid address ID format",
			zap.String("id", id),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	var addressInput model.AddressDTO
	if err := json.NewDecoder(r.Body).Decode(&addressInput); err != nil {
		h.Logger.Error("Failed to decode JSON body",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.Service.AddressService.UpdateAdress(addressId, user.ID, addressInput); err != nil {
		h.Logger.Error("Failed to update address",
			zap.Int("addressId", addressId),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to update address")
		return
	}

	h.Logger.Info("Address successfully updated",
		zap.Int("addressId", addressId),
		zap.String("userId", user.ID),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendSuccess(w, nil, "Address updated successfully")
}

func (h *AddressHandler) SetDefaultAddressHandler(w http.ResponseWriter, r *http.Request) {
	handlerName := "AddressHandler"
	functionName := "SetDefaultAddressHandler"

	if r.Method != http.MethodPatch {
		h.Logger.Error("Invalid method",
			zap.String("method", r.Method),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only PATCH methods are allowed")
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

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Invalid address ID format",
			zap.String("id", id),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	var input struct {
		SetAsDefault bool `json:"set_as_default"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Error("Failed to decode JSON body",
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	if err := h.Service.AddressService.SetDeafultAddress(addressId, user.ID, input.SetAsDefault); err != nil {
		h.Logger.Error("Failed to set default address",
			zap.Int("addressId", addressId),
			zap.Bool("setAsDefault", input.SetAsDefault),
			zap.Error(err),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
		)
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to update address")
		return
	}

	message := "Address successfully set as default"
	if !input.SetAsDefault {
		message = "Address successfully set as not default"
	}
	h.Logger.Info("Default address status updated",
		zap.Int("addressId", addressId),
		zap.Bool("setAsDefault", input.SetAsDefault),
		zap.String("handler", handlerName),
		zap.String("function", functionName),
	)
	JsonResponse.SendSuccess(w, nil, message)
}

func (h *AddressHandler) DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "AddressHandler"
	const functionName = "DeleteAddressHandler"

	if r.Method != http.MethodDelete {
		h.Logger.Error("Invalid request method",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed")
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in the URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID",
			zap.String("id", id),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	err = h.Service.AddressService.RemoveAddress(addressId, user.ID)
	if err != nil {
		h.Logger.Error("Failed to remove address",
			zap.Int("addressId", addressId),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Error(err))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to delete address")
		return
	}

	JsonResponse.SendSuccess(w, nil, "Success delete address")
}

func (h *AddressHandler) GetAllAddressesHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "AddressHandler"
	const functionName = "GetAllAddressesHandler"

	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == "" {
		h.Logger.Error("User ID not found in context",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
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

	userAddress, pagination, err := h.Service.AddressService.GetAllAddresses(user.ID, paginationInput)
	if err != nil {
		h.Logger.Error("Failed to get all addresses",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Error(err))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get all addresses")
		return
	}

	JsonResponse.SendPaginatedResponse(w, userAddress, pagination.Page, pagination.PerPage, pagination.CountData, pagination.CountData/pagination.PerPage, "User address successfully retrieved")
}

func (h *AddressHandler) GetAddressByIdHandler(w http.ResponseWriter, r *http.Request) {
	const handlerName = "AddressHandler"
	const functionName = "GetAddressByIdHandler"

	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method",
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in the URL",
			zap.String("handler", handlerName),
			zap.String("function", functionName))
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID",
			zap.String("id", id),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	address, err := h.Service.AddressService.GetAddressById(addressId)
	if err != nil {
		h.Logger.Error("Failed to get address by ID",
			zap.String("id", id),
			zap.String("handler", handlerName),
			zap.String("function", functionName),
			zap.Error(err))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get address by ID")
		return
	}

	JsonResponse.SendSuccess(w, address, "Address successfully retrieved")
}
