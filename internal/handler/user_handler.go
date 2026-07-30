package handler

import (
	"errors"
	"net/http"

	"github.com/ararext/Go-JWT-Authentication-API/internal/dto"
	"github.com/ararext/Go-JWT-Authentication-API/internal/middleware"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository"
	"github.com/ararext/Go-JWT-Authentication-API/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) Me(c *gin.Context) {
	claims, err := middleware.GetUserClaims(c)
	if err != nil {
		utils.RespondError(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			utils.RespondError(c, http.StatusNotFound, "user not found")
			return
		}
		utils.RespondError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.RespondSuccess(c, http.StatusOK, dto.UserResponse{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	})
}