package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	RedisURL     string
	JWTSecret    string
	Environment  string
}

func Load() *Config {
	// Load .env file
	godotenv.Load()
	
	return &Config{
		Port:         getEnv("USER_SERVICE_PORT", "8001"),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/fintrack?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:    getEnv("JWT_SECRET", "your-secret-key"),
		Environment:  getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}