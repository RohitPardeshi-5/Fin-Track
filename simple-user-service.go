package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"fintrack/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	users = []User{
		{ID: 1, Name: "Test User", Email: "test@example.com"},
	}
	usersMutex sync.RWMutex
	nextUserID = 2
)

func main() {
	r := gin.Default()
	
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "user-service"})
	})

	// Debug endpoint to view all users
	r.GET("/debug/users", func(c *gin.Context) {
		usersMutex.RLock()
		c.JSON(http.StatusOK, gin.H{"users": users, "count": len(users)})
		usersMutex.RUnlock()
	})

	api := r.Group("/api/v1/users")
	
	api.POST("/register", func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		
		usersMutex.Lock()
		// Check for duplicate email
		for _, user := range users {
			if user.Email == req.Email {
				usersMutex.Unlock()
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			}
		}
		
		newUser := User{
			ID:    nextUserID,
			Name:  req.Name,
			Email: req.Email,
		}
		nextUserID++
		users = append(users, newUser)
		usersMutex.Unlock()
		
		c.JSON(http.StatusCreated, gin.H{
			"message": "User created successfully",
			"user":    newUser,
		})
	})

	api.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		usersMutex.RLock()
		for _, user := range users {
			if user.Email == req.Email {
				usersMutex.RUnlock()
				c.JSON(http.StatusOK, gin.H{
					"token": "fake-jwt-token-" + user.Email,
					"user":  user,
				})
				return
			}
		}
		usersMutex.RUnlock()

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	ports := config.LoadPorts()
	addr := ":" + ports.UserService
	fmt.Printf("Starting User Service on %s\n", addr)
	r.Run(addr)
}