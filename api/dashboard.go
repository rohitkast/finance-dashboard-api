package api

import (
	"errors"
	"log"
	"net/http"
	"personal_finance_dashboard/internal/models"
	"personal_finance_dashboard/internal/repository"
	"personal_finance_dashboard/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryParamInput struct {
	Category string `uri:"category" binding:"required,oneof=expense income"`
}

func GetSummary(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var summary *models.SummaryResponse
		summary, err := repository.GetSummary(db, userId)

		if err != nil {
			log.Printf("GetSummary: failed to fetch summary for userId=%d: %v", userId, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Summary fetched successfully",
			"data":    summary,
		})
	}
}

func GetRecent(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var transactions []*models.Transaction

		transactions, err := repository.GetRecent(db, userId)

		if err != nil {
			log.Printf("GetRecent: failed to fetch recent transactions for userId=%d: %v", userId, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Recent transactions fetched successfully",
			"data":    transactions,
		})
	}
}

func GetTransactionsByCategory(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var input CategoryParamInput
		if err := ctx.ShouldBindUri(&input); err != nil {
			log.Printf("GetTransactionsByCategory: invalid category filter for userId=%d: %v", userId, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "category must be either expense or income"})
			return
		}

		transactions, total, err := repository.GetTransactionsByCategory(db, userId, input.Category)
		if err != nil {
			log.Printf("GetTransactionsByCategory: failed for userId=%d category=%s: %v", userId, input.Category, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  input.Category + " transactions fetched successfully",
			"category": input.Category,
			"total":    total,
			"data":     transactions,
		})
	}
}

func GetMonthlyTrends(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		var input CategoryParamInput
		if err := ctx.ShouldBindUri(&input); err != nil {
			log.Printf("GetMonthlyTrends: invalid category for userId=%d: %v", userId, err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "category must be either expense or income"})
			return
		}

		trends, err := repository.GetMonthlyTrends(db, userId, input.Category)
		if err != nil {
			log.Printf("GetMonthlyTrends: failed for userId=%d category=%s: %v", userId, input.Category, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  input.Category + " monthly trends fetched successfully",
			"category": input.Category,
			"data":     trends,
		})
	}
}
