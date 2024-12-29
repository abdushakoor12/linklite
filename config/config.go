package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Host string
		Port string
	}
	DatabaseURL string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	config := &Config{}

	// Server config
	config.Server.Host = getEnv("HOST", "localhost")
	config.Server.Port = getEnv("PORT", "8082")

	// Database config
	config.DatabaseURL = getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/linklite")

	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
