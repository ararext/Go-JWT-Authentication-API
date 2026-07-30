package handler

import (
	"errors"
	"net/http"

	"github.com/ararext/Go-JWT-Authentication-API/internal/dto"
	"github.com/ararext/Go-JWT-Authentication-API/internal/service"
	"github.com/ararext/Go-JWT-Authentication-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req dto.SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Signup(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			utils.RespondError(c, http.StatusConflict, "email already exists")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.RespondSuccess(c, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.RespondError(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(req); err != nil {
		utils.RespondError(c, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Client-side logout: the client discards its tokens.
	// TODO: once TokenBlacklist has a real (e.g. Redis) implementation,
	// blacklist the access token here so it's rejected immediately even
	// before it naturally expires.
	utils.RespondSuccess(c, http.StatusOK, gin.H{"message": "logged out successfully"})
}