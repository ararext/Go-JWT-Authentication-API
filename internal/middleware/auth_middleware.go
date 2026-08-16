package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ararext/Go-JWT-Authentication-API/internal/utils"
	"github.com/gin-gonic/gin"
)

const ContextUserClaimsKey = "userClaims"

func JWTAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			utils.RespondError(c, http.StatusUnauthorized, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondError(c, http.StatusUnauthorized, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := utils.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			if errors.Is(err, utils.ErrExpiredToken) {
				utils.RespondError(c, http.StatusUnauthorized, "token expired")
			} else {
				utils.RespondError(c, http.StatusUnauthorized, "invalid token")
			}
			c.Abort()
			return
		}

		c.Set(ContextUserClaimsKey, claims)
		c.Next()
	}
}
