package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                 string
	MongoURI             string
	DatabaseName         string
	JWTSecret            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on system env vars")
	}

	accessDur, err := time.ParseDuration(getEnv("ACCESS_TOKEN_DURATION", "15m"))
	if err != nil {
		log.Fatalf("invalid ACCESS_TOKEN_DURATION: %v", err)
	}

	refreshDur, err := time.ParseDuration(getEnv("REFRESH_TOKEN_DURATION", "168h"))
	if err != nil {
		log.Fatalf("invalid REFRESH_TOKEN_DURATION: %v", err)
	}

	cfg := &Config{
		Port:                 getEnv("PORT", "8080"),
		MongoURI:             getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName:         getEnv("DATABASE_NAME", "jwt_auth"),
		JWTSecret:            getEnv("JWT_SECRET", ""),
		AccessTokenDuration:  accessDur,
		RefreshTokenDuration: refreshDur,
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}