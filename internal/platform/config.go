package platform

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
)

type Config struct {
	DBURL     string
	Port      string
	JWTSecret string
}

// LoadConfig reads a .env file (if present) and returns a Config struct
func LoadConfig() (*Config, error) {
	// 1. If a .env file exists, load it – this is useful for local dev
	if err := godotenv.Load(); err != nil {
		// It's ok if the file is missing – env vars may be set by the container
		log.Println("⚠️  .env file not found, using OS env vars")
	}

	cfg := &Config{
		DBURL:     os.Getenv("DB_URL"),
		Port:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	// Basic validation
	if cfg.DBURL == "" {
		return nil, fmt.Errorf("DB_URL is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080" //fallback
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return cfg, nil
}
