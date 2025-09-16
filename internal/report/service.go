package report

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"fintrack/internal/common"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service struct {
	db     *gorm.DB
	redis  *redis.Client
	logger *zap.Logger
}

type ReportData struct {
	TotalExpenses float64            `json:"total_expenses"`
	ExpenseCount  int64              `json:"expense_count"`
	Categories    map[string]float64 `json:"categories"`
	Period        string             `json:"period"`
}

func NewService(db *gorm.DB, redis *redis.Client, logger *zap.Logger) *Service {
	return &Service{
		db:     db,
		redis:  redis,
		logger: logger,
	}
}

func (s *Service) GenerateMonthlyReport(userID uint, year int, month int) (*common.Report, error) {
	period := fmt.Sprintf("%d-%02d", year, month)
	
	// Check if report already exists
	var existingReport common.Report
	if err := s.db.Where("user_id = ? AND type = ? AND period = ?", userID, "monthly", period).First(&existingReport).Error; err == nil {
		return &existingReport, nil
	}

	// Generate report in background
	reportChan := make(chan *ReportData, 1)
	errorChan := make(chan error, 1)

	go s.generateReportData(userID, year, month, "monthly", reportChan, errorChan)

	select {
	case reportData := <-reportChan:
		dataJSON, _ := json.Marshal(reportData)
		
		report := common.Report{
			UserID: userID,
			Type:   "monthly",
			Period: period,
			Data:   string(dataJSON),
		}

		if err := s.db.Create(&report).Error; err != nil {
			return nil, err
		}

		// Publish notification
		s.publishNotification(userID, "monthly_report_generated", period)
		
		return &report, nil
	case err := <-errorChan:
		return nil, err
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("report generation timeout")
	}
}

func (s *Service) generateReportData(userID uint, year, month int, reportType string, reportChan chan<- *ReportData, errorChan chan<- error) {
	var expenses []common.Expense
	var total float64
	var count int64
	categories := make(map[string]float64)

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	query := s.db.Where("user_id = ? AND date >= ? AND date < ?", userID, startDate, endDate)
	
	if err := query.Find(&expenses).Error; err != nil {
		errorChan <- err
		return
	}

	if err := query.Count(&count).Error; err != nil {
		errorChan <- err
		return
	}

	for _, expense := range expenses {
		total += expense.Amount
		categories[expense.Category] += expense.Amount
	}

	reportData := &ReportData{
		TotalExpenses: total,
		ExpenseCount:  count,
		Categories:    categories,
		Period:        fmt.Sprintf("%d-%02d", year, month),
	}

	reportChan <- reportData
}

func (s *Service) publishNotification(userID uint, event, data string) {
	notification := map[string]interface{}{
		"user_id": userID,
		"event":   event,
		"data":    data,
		"time":    time.Now(),
	}

	notificationJSON, _ := json.Marshal(notification)
	s.redis.Publish(context.Background(), "notifications", notificationJSON)
}

func (s *Service) GetReports(userID uint, reportType string) ([]common.Report, error) {
	var reports []common.Report
	err := s.db.Where("user_id = ? AND type = ?", userID, reportType).
		Order("created_at DESC").
		Find(&reports).Error
	return reports, err
}