package handler

import (
	"net/http"

	"github.com/soundmarket/backend/internal/domain"
	httprequest "github.com/soundmarket/backend/internal/http/request"
	"github.com/soundmarket/backend/internal/http/middleware"
	"github.com/soundmarket/backend/internal/http/response"
	"github.com/soundmarket/backend/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

type authRequest struct {
	Email    string      `json:"email"`
	Password string      `json:"password"`
	Role     domain.Role `json:"role"`
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.service.Register(req.Email, req.Password, req.Role)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, result)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := httprequest.DecodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid json")
		return
	}
	result, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	user := middleware.CurrentUser(r)
	result, err := h.service.Me(user.ID)
	if err != nil {
		response.FromError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, result)
}
