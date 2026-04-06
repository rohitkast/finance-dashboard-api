package repository

import (
	"context"
	"log"
	"personal_finance_dashboard/internal/models"
	"time"

	"gorm.io/gorm"
)

func GetSummary(db *gorm.DB, userId uint) (*models.SummaryResponse, error) {
	// userid based
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// an address to an actual empty struct instead of a pointer to a nil struct
	var summary *models.SummaryResponse
	if err := db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select(`
		COALESCE(SUM(CASE WHEN category='income' THEN amount ELSE 0 END), 0) AS income,
		COALESCE(SUM(CASE WHEN category='expense' THEN amount ELSE 0 END), 0) AS expense
		`).
		Where("user_id=?", userId).
		Scan(&summary).
		Error; err != nil {
		log.Printf("repository.GetSummary: failed for userId=%d: %v", userId, err)
		return nil, err
	}

	summary.Balance = summary.Income - summary.Expense
	return summary, nil
}

// get last 10 transactions
func GetRecent(db *gorm.DB, userId uint) ([]*models.Transaction, error) {
	// userid based
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transactions []*models.Transaction
	if err := db.WithContext(ctx).
		Find(&transactions).
		Where("user_id=?", userId).
		Limit(10).
		Error; err != nil {
		log.Printf("repository.GetRecent: failed for userId=%d: %v", userId, err)
		return nil, err
	}

	return transactions, nil
}

func GetTransactionsByCategory(db *gorm.DB, userId uint, category string) ([]*models.Transaction, float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transactions []*models.Transaction
	if err := db.WithContext(ctx).
		Where("user_id = ? AND category = ?", userId, category).
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {
		log.Printf("repository.GetTransactionsByCategory: failed to fetch transactions for userId=%d category=%s: %v", userId, category, err)
		return nil, 0, err
	}

	// instead of total being on every row we are calculating it seperately
	var total float64
	if err := db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("user_id = ? AND category = ?", userId, category).
		Scan(&total).Error; err != nil {
		log.Printf("repository.GetTransactionsByCategory: failed to calculate total for userId=%d category=%s: %v", userId, category, err)
		return nil, 0, err
	}

	return transactions, total, nil
}

func GetMonthlyTrends(db *gorm.DB, userId uint, category string) ([]*models.MonthlyTrendResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var trends []*models.MonthlyTrendResponse
	if err := db.WithContext(ctx).
		Model(&models.Transaction{}).
		Select("EXTRACT(MONTH FROM created_at)::int AS month, COALESCE(SUM(amount), 0) AS total").
		Where("user_id = ? AND category = ?", userId, category).
		Group("EXTRACT(MONTH FROM created_at)").
		Order("month ASC").
		Scan(&trends).Error; err != nil {
		log.Printf("repository.GetMonthlyTrends: failed for userId=%d category=%s: %v", userId, category, err)
		return nil, err
	}

	return trends, nil
}
