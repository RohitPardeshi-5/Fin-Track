package expense

import (
	"testing"
	"fintrack/internal/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&common.Expense{})
	return db
}

func TestExpenseService_CreateExpense(t *testing.T) {
	db := setupTestDB()
	service := NewService(db)

	req := common.ExpenseRequest{
		Amount:      100.50,
		Description: "Test expense",
		Category:    "Food",
		Date:        "2024-01-15",
	}

	expense, err := service.CreateExpense(1, req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if expense.Amount != req.Amount {
		t.Errorf("Expected amount %f, got %f", req.Amount, expense.Amount)
	}

	if expense.Category != req.Category {
		t.Errorf("Expected category %s, got %s", req.Category, expense.Category)
	}
}

func TestExpenseService_GetExpenses(t *testing.T) {
	db := setupTestDB()
	service := NewService(db)

	// Create test expense
	req := common.ExpenseRequest{
		Amount:      100.50,
		Description: "Test expense",
		Category:    "Food",
		Date:        "2024-01-15",
	}
	service.CreateExpense(1, req)

	expenses, err := service.GetExpenses(1, 10, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(expenses) != 1 {
		t.Errorf("Expected 1 expense, got %d", len(expenses))
	}
}