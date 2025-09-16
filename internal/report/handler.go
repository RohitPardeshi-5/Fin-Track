package report

import (
	"net/http"
	"strconv"
	"time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *Service
	logger  *zap.Logger
}

func NewHandler(service *Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) GetMonthlyReport(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	monthStr := c.DefaultQuery("month", strconv.Itoa(int(time.Now().Month())))
	
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
		return
	}
	
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
		return
	}

	report, err := h.service.GenerateMonthlyReport(userID, year, month)
	if err != nil {
		h.logger.Error("Failed to generate monthly report", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetReports(c *gin.Context) {
	userID := c.GetUint("user_id")
	reportType := c.DefaultQuery("type", "monthly")

	reports, err := h.service.GetReports(userID, reportType)
	if err != nil {
		h.logger.Error("Failed to get reports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/monthly", h.GetMonthlyReport)
	router.GET("", h.GetReports)
}