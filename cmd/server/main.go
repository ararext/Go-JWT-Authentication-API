package main

import (
	"net/http"

	"github.com/ararext/Go-JWT-Authentication-API/internal/config"
	"github.com/ararext/Go-JWT-Authentication-API/internal/database"
	"github.com/ararext/Go-JWT-Authentication-API/internal/handler"
	"github.com/ararext/Go-JWT-Authentication-API/internal/logger"
	"github.com/ararext/Go-JWT-Authentication-API/internal/repository"
	"github.com/ararext/Go-JWT-Authentication-API/internal/routes"
	"github.com/ararext/Go-JWT-Authentication-API/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New()
	defer log.Sync()

	log.Info("starting server",
		zap.String("port", cfg.Port),
		zap.String("database", cfg.DatabaseName),
	)

	db, err := database.Connect(cfg.MongoURI, cfg.DatabaseName, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	// Dependency chain: repository -> service -> handler
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(
		userRepo,
		cfg.JWTSecret,
		cfg.AccessTokenDuration,
		cfg.RefreshTokenDuration,
	)
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userRepo)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routes.RegisterRoutes(router, authHandler, userHandler, cfg.JWTSecret)

	addr := ":" + cfg.Port
	log.Info("server listening", zap.String("address", addr))

	if err := router.Run(addr); err != nil {
		log.Fatal("server failed to start", zap.Error(err))
	}
}