package repository

import (
	"context"
	"log"
	"personal_finance_dashboard/internal/models"
	"time"

	"gorm.io/gorm"
)

func CreateTransaction(db *gorm.DB, transaction *models.Transaction) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.WithContext(ctx).Create(transaction).Error; err != nil {
		log.Printf("repository.CreateTransaction: failed for userId=%d: %v", transaction.UserID, err)
		return nil, err
	}
	return transaction, nil
}

func GetAllTransactions(db *gorm.DB, userID uint) ([]*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transactions []*models.Transaction
	if err := db.WithContext(ctx).
		Where("user_id=?", userID).
		Find(&transactions).
		Error; err != nil {
		log.Printf("repository.GetAllTransactions: failed: %v", err)
		return nil, err
	}
	return transactions, nil
}

func GetTransactionById(db *gorm.DB, id uint64, uid uint) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction *models.Transaction
	if err := db.WithContext(ctx).
		Where("id=? AND user_id=?", id, uid).
		First(&transaction).
		Error; err != nil {
		log.Printf("repository.GetTransactionById: failed for id=%d: %v", id, err)
		return nil, err
	}
	return transaction, nil
}

func UpdateTransaction(db *gorm.DB, updates *models.Transaction, id uint64) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction *models.Transaction
	if err := db.WithContext(ctx).
		Where("id=? AND user_id=?", id, updates.UserID).
		First(&transaction).Error; err != nil {
		log.Printf("repository.UpdateTransaction: failed lookup for id=%d: %v", id, err)
		return nil, err // Returns gorm.ErrRecordNotFound if not found
	}

	if err := db.WithContext(ctx).Model(&transaction).Updates(updates).Error; err != nil {
		log.Printf("repository.UpdateTransaction: failed update for id=%d: %v", id, err)
		return nil, err
	}
	return transaction, nil
}

// this is soft delete only
func DeleteTransaction(db *gorm.DB, id uint64, uid uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transaction, err := GetTransactionById(db, id, uid)
	if err != nil {
		log.Printf("repository.DeleteTransaction: failed lookup for id=%d: %v", id, err)
		return err // Returns gorm.ErrRecordNotFound if not found
	}

	if err := db.WithContext(ctx).Delete(&transaction).Error; err != nil {
		log.Printf("repository.DeleteTransaction: failed delete for id=%d: %v", id, err)
		return err
	}

	return nil
}

func GetFilteredTransactions(db *gorm.DB, uid uint, asc bool, dsc bool, month string, year string, from string, to string, exp bool, inc bool) ([]*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := db.WithContext(ctx).Model(&models.Transaction{}).Where("user_id = ?", uid)

	if month != "" {
		query = query.Where("EXTRACT(MONTH FROM created_at) = ?", month)
	}

	if year != "" {
		query = query.Where("EXTRACT(YEAR FROM created_at) = ?", year)
	}

	if from != "" {
		query = query.Where("DATE(created_at) >= ?", from)
	}

	if to != "" {
		query = query.Where("DATE(created_at) <= ?", to)
	}

	if exp && !inc {
		query = query.Where("category = ?", "expense")
	} else if inc && !exp {
		query = query.Where("category = ?", "income")
	}

	if asc && !dsc {
		query = query.Order("amount ASC")
	} else if dsc && !asc {
		query = query.Order("amount DESC")
	} else {
		query = query.Order("created_at DESC")
	}

	var transactions []*models.Transaction
	if err := query.Find(&transactions).Error; err != nil {
		log.Printf("repository.GetFilteredTransactions: failed for userId=%d: %v", uid, err)
		return nil, err
	}

	return transactions, nil
}
