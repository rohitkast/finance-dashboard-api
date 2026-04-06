package repository

import (
	"context"
	"log"
	"personal_finance_dashboard/internal/models"
	"time"

	"gorm.io/gorm"
)

func AdminDeleteTransaction(db *gorm.DB, id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var transaction *models.Transaction
	if err := db.WithContext(ctx).First(&transaction, id).Error; err != nil {
		log.Printf("repository.DeleteTransaction: failed lookup for id=%d: %v", id, err)
		return err
	}

	if err := db.WithContext(ctx).Delete(&transaction).Error; err != nil {
		log.Printf("repository.DeleteTransaction: failed delete for id=%d: %v", id, err)
		return err
	}

	return nil
}
