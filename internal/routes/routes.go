package routes

import (
	"github.com/ararext/Go-JWT-Authentication-API/internal/handler"
	"github.com/ararext/Go-JWT-Authentication-API/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authHandler *handler.AuthHandler, userHandler *handler.UserHandler, jwtSecret string) {
	v1 := router.Group("/api/v1")

	auth := v1.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}

	users := v1.Group("/users")
	users.Use(middleware.JWTAuth(jwtSecret))
	{
		users.GET("/me", userHandler.Me)
	}
}