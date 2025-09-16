package expense

import (
	"time"
	"fintrack/internal/common"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateExpense(userID uint, req common.ExpenseRequest) (*common.Expense, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	expense := common.Expense{
		UserID:      userID,
		Amount:      req.Amount,
		Description: req.Description,
		Category:    req.Category,
		Date:        date,
	}

	if err := s.db.Create(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (s *Service) GetExpenses(userID uint, limit, offset int) ([]common.Expense, error) {
	var expenses []common.Expense
	err := s.db.Where("user_id = ?", userID).
		Order("date DESC").
		Limit(limit).
		Offset(offset).
		Find(&expenses).Error
	return expenses, err
}

func (s *Service) UpdateExpense(userID, expenseID uint, req common.ExpenseRequest) (*common.Expense, error) {
	var expense common.Expense
	if err := s.db.Where("id = ? AND user_id = ?", expenseID, userID).First(&expense).Error; err != nil {
		return nil, err
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, err
	}

	expense.Amount = req.Amount
	expense.Description = req.Description
	expense.Category = req.Category
	expense.Date = date

	if err := s.db.Save(&expense).Error; err != nil {
		return nil, err
	}

	return &expense, nil
}

func (s *Service) DeleteExpense(userID, expenseID uint) error {
	return s.db.Where("id = ? AND user_id = ?", expenseID, userID).Delete(&common.Expense{}).Error
}