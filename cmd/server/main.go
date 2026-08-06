package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ararext/Go-JWT-Authentication-API/internal/config"
	"github.com/ararext/Go-JWT-Authentication-API/internal/database"
	"github.com/ararext/Go-JWT-Authentication-API/internal/handler"
	"github.com/ararext/Go-JWT-Authentication-API/internal/logger"
	"github.com/ararext/Go-JWT-Authentication-API/internal/middleware"
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

	// Security middleware — order matters: CORS/headers first, then rate limiting,
	// so even rejected requests get consistent headers.
	router.Use(middleware.CORS())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.RateLimit(5, 10)) // 5 req/sec per IP, burst of 10

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	routes.RegisterRoutes(router, authHandler, userHandler, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Run the server in a goroutine so it doesn't block shutdown handling.
	go func() {
		log.Info("server listening", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed to start", zap.Error(err))
		}
	}()

	// Block until an interrupt or terminate signal is received.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("forced server shutdown", zap.Error(err))
	}

	log.Info("server exited cleanly")
}