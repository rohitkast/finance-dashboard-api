package api

import (
	"errors"
	"log"
	"net/http"
	"personal_finance_dashboard/internal/models"
	"personal_finance_dashboard/internal/repository"
	"personal_finance_dashboard/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateTransactionInput struct {
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
	Category    string  `json:"category" binding:"required,oneof=income expense"`
}

type FilterTransactionInput struct {
	Asc   bool   `form:"asc"`
	Dsc   bool   `form:"dsc"`
	Month string `form:"month"`
	Year  string `form:"year"`
	From  string `form:"from"`
	To    string `form:"to"`
	Exp   bool   `form:"exp"`
	Inc   bool   `form:"inc"`
}

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var input CreateTransactionInput
		// validate input
		if err := ctx.ShouldBindJSON(&input); err != nil {
			log.Printf("CreateTransaction: invalid request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		create := &models.Transaction{
			UserID:      userId,
			Amount:      input.Amount,
			Description: input.Description,
			Category:    input.Category,
		}

		transaction, err := repository.CreateTransaction(db, create)
		if err != nil {
			log.Printf("CreateTransaction: failed to create transaction for userId=%d: %v", create.UserID, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"success":             true,
			"message":             "Transaction created successfully",
			"transaction_details": transaction,
		})
	}

}

func GetAllTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var transactions []*models.Transaction

		transactions, err := repository.GetAllTransactions(db, userId)

		if err != nil {
			log.Printf("GetAllTransactions: failed to fetch transactions: %v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transactions fetched successfully",
			"data":    transactions,
		})
	}

}

func GetTransactionById(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idstr := ctx.Param("id")
		id, err := strconv.ParseUint(idstr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("GetTransactionById: invalid id=%q", idstr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		var transaction *models.Transaction

		transaction, err = repository.GetTransactionById(db, id, userId)
		if err != nil {
			log.Printf("GetTransactionById: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction fetched successfully",
			"data":    transaction,
		})
	}
}

func UpdateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idstr := ctx.Param("id")
		id, err := strconv.ParseUint(idstr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("UpdateTransaction: invalid id=%q", idstr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		var input CreateTransactionInput
		// validate input
		if err := ctx.ShouldBindJSON(&input); err != nil {
			log.Printf("UpdateTransaction: invalid request body for id=%d: %v", id, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		update := &models.Transaction{
			UserID:      userId,
			Amount:      input.Amount,
			Description: input.Description,
			Category:    input.Category,
		}

		var transaction *models.Transaction

		transaction, err = repository.UpdateTransaction(db, update, id)
		if err != nil {
			log.Printf("UpdateTransaction: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction updated successfully",
			"data":    transaction,
		})
	}
}

func DeleteTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idstr := ctx.Param("id")
		id, err := strconv.ParseUint(idstr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("DeleteTransaction: invalid id=%q", idstr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		err = repository.DeleteTransaction(db, id, userId)
		if err != nil {
			log.Printf("DeleteTransaction: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Transaction deleted successfully",
			"success": true,
		})
	}
}

func GetFilteredTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var filters FilterTransactionInput
		if err := ctx.ShouldBindQuery(&filters); err != nil {
			log.Printf("GetFilteredTransactions: invalid query params: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		if filters.From != "" {
			if _, err := time.Parse("2006-01-02", filters.From); err != nil {
				log.Printf("GetFilteredTransactions: invalid from date=%q", filters.From)
				ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "from must be in YYYY-MM-DD format"})
				return
			}
		}

		if filters.To != "" {
			if _, err := time.Parse("2006-01-02", filters.To); err != nil {
				log.Printf("GetFilteredTransactions: invalid to date=%q", filters.To)
				ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "to must be in YYYY-MM-DD format"})
				return
			}
		}

		transactions, err := repository.GetFilteredTransactions(db, userId, filters.Asc, filters.Dsc, filters.Month, filters.Year, filters.From, filters.To, filters.Exp, filters.Inc)
		if err != nil {
			log.Printf("GetFilteredTransactions: failed to fetch filtered transactions for userId=%d: %v", userId, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transactions fetched successfully",
			"data":    transactions,
		})
	}
}
