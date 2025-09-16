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

type Report struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Period    string `json:"period"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

var (
	reports = []Report{
		{ID: 1, Type: "Monthly", Period: "January 2024", Status: "completed", CreatedAt: "2024-01-31"},
	}
	reportsMutex sync.RWMutex
	nextReportID = 2
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
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "report-service"})
	})

	api := r.Group("/api/v1")
	reportsAPI := api.Group("/reports")
	
	reportsAPI.GET("", func(c *gin.Context) {
		reportsMutex.RLock()
		defer reportsMutex.RUnlock()
		c.JSON(http.StatusOK, gin.H{"reports": reports})
	})

	reportsAPI.GET("/monthly", func(c *gin.Context) {
		reportsMutex.Lock()
		newReport := Report{
			ID:        nextReportID,
			Type:      "Monthly",
			Period:    time.Now().Format("January 2006"),
			Status:    "completed",
			CreatedAt: time.Now().Format("2006-01-02"),
		}
		nextReportID++
		reports = append(reports, newReport)
		reportsMutex.Unlock()
		
		c.JSON(http.StatusOK, gin.H{
			"message": "Monthly report generated",
			"report":  newReport,
		})
	})

	ports := config.LoadPorts()
	addr := ":" + ports.ReportService
	fmt.Printf("Starting Report Service on %s\n", addr)
	r.Run(addr)
}