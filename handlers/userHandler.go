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

	userID, err := h.Service.UserService.CreateUser(userInput)
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

	user, err := h.Service.UserService.Login(userInput)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := util.GenerateToken(user.ID, h.Config)
	if err != nil {
		h.Logger.Error(err.Error())
		JsonResponse.SendError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// save user id in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
	})

	JsonResponse.SendSuccess(w, map[string]interface{}{
		"token": token,
	}, "Login successful")
}
