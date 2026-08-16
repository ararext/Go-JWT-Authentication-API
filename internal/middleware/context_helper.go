package middleware

import (
	"errors"

	"github.com/ararext/Go-JWT-Authentication-API/internal/utils"
	"github.com/gin-gonic/gin"
)

var ErrClaimsNotFound = errors.New("claims not found in context")

func GetUserClaims(c *gin.Context) (*utils.Claims, error) {
	value, exists := c.Get(ContextUserClaimsKey)
	if !exists {
		return nil, ErrClaimsNotFound
	}

	claims, ok := value.(*utils.Claims)
	if !ok {
		return nil, ErrClaimsNotFound
	}

	return claims, nil
}
