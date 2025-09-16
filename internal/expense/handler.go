package expense

import (
	"net/http"
	"strconv"
	"fintrack/internal/common"
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

func (h *Handler) CreateExpense(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req common.ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, err := h.service.CreateExpense(userID, req)
	if err != nil {
		h.logger.Error("Failed to create expense", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, expense)
}

func (h *Handler) GetExpenses(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	expenses, err := h.service.GetExpenses(userID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to get expenses", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func (h *Handler) UpdateExpense(c *gin.Context) {
	userID := c.GetUint("user_id")
	expenseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	var req common.ExpenseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	expense, err := h.service.UpdateExpense(userID, uint(expenseID), req)
	if err != nil {
		h.logger.Error("Failed to update expense", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, expense)
}

func (h *Handler) DeleteExpense(c *gin.Context) {
	userID := c.GetUint("user_id")
	expenseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	if err := h.service.DeleteExpense(userID, uint(expenseID)); err != nil {
		h.logger.Error("Failed to delete expense", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted successfully"})
}

func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	router.POST("", h.CreateExpense)
	router.GET("", h.GetExpenses)
	router.PUT("/:id", h.UpdateExpense)
	router.DELETE("/:id", h.DeleteExpense)
}