package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/helper"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/model"
	"github.com/Safiramdhn/project-app-ecommerce-golang-safira/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	UserService service.UserService
	Logger      *zap.Logger
}

var JsonResponse = helper.JSONResponse{}

func NewUserHandler(userService service.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{UserService: userService, Logger: log}
}

func (h *UserHandler) RegisterHanlder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error(errMessage)
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
	}

	var userInput model.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID, err := h.UserService.CreateUser(userInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	JsonResponse.SendCreated(w, userID, "User created successfully")
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errMessage := fmt.Sprintf("Invalid method %s", r.Method)
		h.Logger.Error(errMessage)
		JsonResponse.SendError(w, http.StatusBadRequest, errMessage)
	}

	var userInput model.UserDTO
	err := json.NewDecoder(r.Body).Decode(&userInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user, err := h.UserService.Login(userInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}
	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	JsonResponse.SendSuccess(w, map[string]interface{}{
		"token": token,
	}, "Login successful")
}
