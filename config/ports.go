package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Ports struct {
	UserService     string
	ExpenseService  string
	ReportService   string
	WebFrontend     string
}

func LoadPorts() *Ports {
	godotenv.Load()
	
	return &Ports{
		UserService:     getEnv("USER_SERVICE_PORT", "8001"),
		ExpenseService:  getEnv("EXPENSE_SERVICE_PORT", "8002"),
		ReportService:   getEnv("REPORT_SERVICE_PORT", "8003"),
		WebFrontend:     getEnv("WEB_FRONTEND_PORT", "8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}