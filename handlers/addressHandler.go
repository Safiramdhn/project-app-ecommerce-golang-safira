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

type AddressHandler struct {
	Service service.MainService
	Logger  *zap.Logger
}

func NewAddressHandler(service service.MainService, log *zap.Logger) AddressHandler {
	return AddressHandler{Service: service, Logger: log}
}

func (h *AddressHandler) AddAddressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "RegisterHandler"))
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

	var addressInput model.AddressDTO
	err := json.NewDecoder(r.Body).Decode(&addressInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.AddressService.AddAddress(user.ID, addressInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	JsonResponse.SendCreated(w, nil, "Address successfully created")
}

func (h *AddressHandler) UpdateAddressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler",
			"UpdateAddressHandler"), zap.String("function", "UpdateAddressHandler"))
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

	var addressInput model.AddressDTO
	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&addressInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.Service.AddressService.UpdateAdress(addressId, user.ID, addressInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler",
			"Address"), zap.String("function", "AddAddressHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	JsonResponse.SendSuccess(w, nil, "Address updated successfully")
}

func (h *AddressHandler) SetDefaultAddressHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is PATCH
	if r.Method != http.MethodPatch {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "Address"), zap.String("function", "SetDefaultAddressHandler"))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, errMessage)
		return
	}

	// Retrieve user context
	ctxValue := r.Context().Value(middleware.UserClaimsContextKey)
	if ctxValue == nil {
		h.Logger.Error("User ID not found in context")
		JsonResponse.SendError(w, http.StatusBadRequest, "User ID missing or invalid")
		return
	}

	user, ok := ctxValue.(model.User)
	if !ok {
		h.Logger.Error("Failed to cast user from context")
		JsonResponse.SendError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Parse and validate URL parameters
	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	var setAsDefault bool
	err = json.NewDecoder(r.Body).Decode(&setAsDefault)
	if err != nil {
		h.Logger.Error("Failed to decode JSON body", zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "invalid payload", err)
		return
	}

	// Call the service to set the default address
	err = h.Service.AddressService.SetDeafultAddress(addressId, user.ID, setAsDefault)
	if err != nil {
		h.Logger.Error("Failed to set default address", zap.Int("addressId", addressId), zap.Bool("setAsDefault", setAsDefault), zap.Error(err))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to update address")
		return
	}

	// Send success response
	message := "Address successfully set as default"
	if !setAsDefault {
		message = "Address successfully set as not default"
	}
	JsonResponse.SendSuccess(w, nil, message)
}

func (h *AddressHandler) DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed")
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

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	err = h.Service.AddressService.RemoveAddress(addressId, user.ID)
	if err != nil {
		h.Logger.Error("Failed to remove address", zap.Int("addressId", addressId), zap.String("handler",
			"Address"), zap.String("function", "DeleteAddressHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to delete address", err)
		return
	}

	JsonResponse.SendSuccess(w, nil, "Success delete address")
}

func (h *AddressHandler) GetAllAddressesHandler(w http.ResponseWriter, r *http.Request) {
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
		h.Logger.Error("Failed to get all addresses", zap.String("handler", "Address"), zap.String("function", "GetAllAddressesHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get all addresses", err)
		return
	}

	if pagination.CountData/pagination.PerPage > 0 {
		TotalPage = pagination.CountData / pagination.PerPage
	}
	JsonResponse.SendPaginatedResponse(w, userAddress, pagination.Page, pagination.PerPage, pagination.CountData, TotalPage, "User address successfully retrieved")
}

func (h *AddressHandler) GetAddressByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.Logger.Error("Invalid request method", zap.String("method", r.Method))
		JsonResponse.SendError(w, http.StatusMethodNotAllowed, "Only GET methods are allowed")
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		h.Logger.Error("Address ID is missing in the URL")
		JsonResponse.SendError(w, http.StatusBadRequest, "Address ID is required")
		return
	}

	addressId, err := strconv.Atoi(id)
	if err != nil {
		h.Logger.Error("Failed to parse address ID", zap.String("id", id), zap.Error(err))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid address ID format")
		return
	}

	address, err := h.Service.AddressService.GetAddressById(addressId)
	if err != nil {
		h.Logger.Error("Failed to get address by ID", zap.String("id", id), zap.String("handler", "address"), zap.String("functions", "GetAddressByIdHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to get address by ID")
		return
	}

	JsonResponse.SendSuccess(w, address, "Address successfully retrieved")
}
