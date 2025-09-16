package main

import (
	"fmt"
	"log"
	"net/http"
	"fintrack/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.LoadHTMLFiles(
		"web/templates/auth-page.html",
		"web/templates/modern-dashboard.html",
		"web/templates/modern-expenses.html",
		"web/templates/modern-reports.html",
		"web/templates/simple-ai-chat.html",
	)
	r.Static("/static", "./web/static")
	
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "web-frontend"})
	})
	
	// Handle favicon.ico to prevent 404 errors
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	
	r.GET("/", handleLogin)
	r.GET("/login", handleLogin)
	r.GET("/register", handleRegister)
	r.GET("/dashboard", handleDashboard)
	r.GET("/expenses", handleExpenses)
	r.GET("/reports", handleReports)
	r.GET("/ai-chat", handleAIChat)
	
	ports := config.LoadPorts()
	addr := ":" + ports.WebFrontend
	fmt.Printf("Starting Web Frontend on %s\n", addr)
	log.Printf("Frontend server starting on %s", addr)
	log.Fatal(r.Run(addr))
}

func handleHome(c *gin.Context) {
	c.HTML(http.StatusOK, "standalone-index.html", gin.H{"title": "FinTrack"})
}

func handleLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "auth-page.html", gin.H{"title": "Authentication"})
}

func handleRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth-page.html", gin.H{"title": "Authentication"})
}



func handleDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "modern-dashboard.html", gin.H{"title": "Dashboard"})
}

func handleExpenses(c *gin.Context) {
	c.HTML(http.StatusOK, "modern-expenses.html", gin.H{"title": "Expenses"})
}

func handleReports(c *gin.Context) {
	c.HTML(http.StatusOK, "modern-reports.html", gin.H{"title": "Reports"})
}

func handleAIChat(c *gin.Context) {
	c.HTML(http.StatusOK, "simple-ai-chat.html", gin.H{"title": "AI Assistant"})
}