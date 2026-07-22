package main

import (
	"net/http"

	"github.com/ararext/Go-JWT-Authentication-API/internal/config"
	"github.com/ararext/Go-JWT-Authentication-API/internal/database"
	"github.com/ararext/Go-JWT-Authentication-API/internal/logger"
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
	_ = db // used properly starting Day 6 once repository/service are wired through routes

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	addr := ":" + cfg.Port
	log.Info("server listening", zap.String("address", addr))

	if err := router.Run(addr); err != nil {
		log.Fatal("server failed to start", zap.Error(err))
	}
}