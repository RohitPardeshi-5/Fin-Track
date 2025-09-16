package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
	"fintrack/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Expense struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	UserID      int     `json:"user_id"`
}

var (
	expenses = []Expense{
		{ID: 1, Description: "Lunch", Amount: 25.50, Category: "Food", Date: "2024-01-15", UserID: 1},
		{ID: 2, Description: "Gas", Amount: 45.00, Category: "Transport", Date: "2024-01-14", UserID: 1},
	}
	expensesMutex sync.RWMutex
	nextExpenseID = 3
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
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "expense-service"})
	})

	api := r.Group("/api/v1")
	expensesAPI := api.Group("/expenses")
	
	expensesAPI.GET("", func(c *gin.Context) {
		expensesMutex.RLock()
		defer expensesMutex.RUnlock()
		c.JSON(http.StatusOK, gin.H{"expenses": expenses})
	})

	expensesAPI.POST("", func(c *gin.Context) {
		var expense Expense
		if err := c.ShouldBindJSON(&expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		
		expensesMutex.Lock()
		expense.ID = nextExpenseID
		nextExpenseID++
		expense.UserID = 1
		expenses = append(expenses, expense)
		expensesMutex.Unlock()
		
		c.JSON(http.StatusCreated, gin.H{"expense": expense})
	})

	expensesAPI.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		
		expensesMutex.Lock()
		defer expensesMutex.Unlock()
		
		for i, expense := range expenses {
			if expense.ID == id {
				if err := c.ShouldBindJSON(&expenses[i]); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
					return
				}
				expenses[i].ID = id
				c.JSON(http.StatusOK, gin.H{"expense": expenses[i]})
				return
			}
		}
		
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
	})

	expensesAPI.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		
		expensesMutex.Lock()
		defer expensesMutex.Unlock()
		
		for i, expense := range expenses {
			if expense.ID == id {
				expenses = append(expenses[:i], expenses[i+1:]...)
				c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
				return
			}
		}
		
		c.JSON(http.StatusNotFound, gin.H{"error": "Expense not found"})
	})

	ports := config.LoadPorts()
	addr := ":" + ports.ExpenseService
	fmt.Printf("Starting Expense Service on %s\n", addr)
	r.Run(addr)
}