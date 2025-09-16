package user

import (
	"testing"
	"fintrack/internal/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&common.User{})
	return db
}

func TestUserService_Register(t *testing.T) {
	db := setupTestDB()
	service := NewService(db, "test-secret")

	req := common.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}

	response, err := service.Register(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.User.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, response.User.Email)
	}

	if response.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestUserService_Login(t *testing.T) {
	db := setupTestDB()
	service := NewService(db, "test-secret")

	// First register a user
	registerReq := common.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	}
	service.Register(registerReq)

	// Then try to login
	loginReq := common.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	response, err := service.Login(loginReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if response.User.Email != loginReq.Email {
		t.Errorf("Expected email %s, got %s", loginReq.Email, response.User.Email)
	}
}