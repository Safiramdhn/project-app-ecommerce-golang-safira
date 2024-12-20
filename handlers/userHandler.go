package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/util"
	"go.uber.org/zap"
)

type UserHandler struct {
	Service service.MainService
	Logger  *zap.Logger
	Config  util.Configuration
}

var JsonResponse = helper.JSONResponse{}

func NewUserHandler(service service.MainService, log *zap.Logger, config util.Configuration) UserHandler {
	return UserHandler{Service: service, Logger: log, Config: config}
}

func (h *UserHandler) RegisterHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "RegisterHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
	}

	var userInput model.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "GetProductByIdHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	errField := helper.EmailOrPhoneValidator(userInput.EmailOrPhoneNumber)
	if errField.Message != "" {
		// Log the specific validation error
		h.Logger.Error("Validation error",
			zap.String("field", errField.Field),
			zap.String("message", errField.Message),
			zap.String("method", r.Method),
			zap.String("handler", "User"),
			zap.String("function", "GetProductByIdHandler"),
		)
		// Send detailed error response
		JsonResponse.SendError(w, http.StatusBadRequest, errField.Message)
		return
	}

	errField = helper.PasswordValidator(userInput.Password)
	if errField.Message != "" {
		// Log the specific validation error
		h.Logger.Error("Validation error",
			zap.String("field", errField.Field),
			zap.String("message", errField.Message),
			zap.String("method", r.Method),
			zap.String("handler", "User"),
			zap.String("function", "GetProductByIdHandler"),
		)

		// Send detailed error response
		JsonResponse.SendError(w, http.StatusBadRequest, errField.Message)
		return
	}

	userID, err := h.Service.UserService.CreateUser(userInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "RegisterHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	JsonResponse.SendCreated(w, userID, "User created successfully")
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error("Invalid method", zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "LoginHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
	}

	var userInput model.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "LoginHandler"))
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	errField := helper.EmailOrPhoneValidator(userInput.EmailOrPhoneNumber)
	if errField.Message != "" {
		// Log the specific validation error
		h.Logger.Error("Validation error",
			zap.String("field", errField.Field),
			zap.String("message", errField.Message),
			zap.String("method", r.Method),
			zap.String("handler", "User"),
			zap.String("function", "LoginHandler"))
		// Send detailed error response
		JsonResponse.SendError(w, http.StatusBadRequest, errField.Message)
		return
	}

	user, err := h.Service.UserService.Login(userInput)
	if err != nil {
		h.Logger.Error("Authentication error",
			zap.String("error", err.Error()),
			zap.String("method", r.Method),
			zap.String("handler", "User"),
			zap.String("function", "LoginHandler"),
		)
		JsonResponse.SendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if (user == model.User{}) {
		h.Logger.Error("User not found",
			zap.String("email_or_phone", userInput.EmailOrPhoneNumber),
			zap.String("method", r.Method),
			zap.String("handler", "User"),
			zap.String("function", "LoginHandler"),
		)
		JsonResponse.SendError(w, http.StatusNotFound, "User not found")
		return
	}

	token, err := util.GenerateToken(user.ID, h.Config)
	if err != nil {
		h.Logger.Error(err.Error(), zap.String("method", r.Method), zap.String("handler", "User"), zap.String("function", "LoginHandler"))
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	JsonResponse.SendSuccess(w, map[string]interface{}{
		"token": token,
	}, "Login successful")
}
