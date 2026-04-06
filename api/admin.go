package api

import (
	"errors"
	"log"
	"net/http"
	"personal_finance_dashboard/internal/repository"
	"personal_finance_dashboard/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminDeleteTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("AdminDeleteTransaction: invalid id=%q", idStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid transaction id"})
			return
		}

		err = repository.AdminDeleteTransaction(db, id)
		if err != nil {
			log.Printf("AdminDeleteTransaction: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Transaction deleted successfully",
		})
	}
}

func AdminGetUserTransactions(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		adminID := ctx.GetUint("user_id")
		if adminID == 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "user id not found in context",
			})
			return
		}

		uidStr := ctx.Param("uid")
		uid, err := strconv.ParseUint(uidStr, 10, 32)
		if err != nil || uid == 0 {
			log.Printf("AdminGetUserTransactions: invalid uid=%q", uidStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		transactions, err := repository.GetAllTransactions(db, uint(uid))
		if err != nil {
			log.Printf("AdminGetUserTransactions: failed for uid=%d: %v", uid, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User transactions fetched successfully",
			"uid":     uid,
			"data":    transactions,
		})
	}
}

func GetAllUsers(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		users, err := repository.GetAllUsers(db)
		if err != nil {
			log.Printf("AdminGetAllUsers: failed to fetch users: %v", err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "All users fetched successfully",
			"count":   len(users),
			"data":    users,
		})
	}
}
