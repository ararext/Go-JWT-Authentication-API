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

// Signup godoc
// @Summary      Register a new user
// @Description  Creates a new user account and returns access + refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.SignupRequest true "Signup payload"
// @Success      201  {object}  dto.AuthResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Router       /api/v1/auth/signup [post]
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

// Login godoc
// @Summary      Authenticate a user
// @Description  Verifies credentials and returns access + refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login payload"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/v1/auth/login [post]
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

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchanges a valid refresh token for a new access + refresh token pair
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RefreshRequest true "Refresh payload"
// @Success      200  {object}  dto.AuthResponse
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/v1/auth/refresh [post]
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

// Logout godoc
// @Summary      Log out
// @Description  Client-side logout (stateless for now — token blacklisting to be added later)
// @Tags         auth
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Client-side logout: the client discards its tokens.
	// TODO: once TokenBlacklist has a real (e.g. Redis) implementation,
	// blacklist the access token here so it's rejected immediately even
	// before it naturally expires.
	utils.RespondSuccess(c, http.StatusOK, gin.H{"message": "logged out successfully"})
}
