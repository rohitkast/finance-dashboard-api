package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// retrieves the authenticated user ID from the context
func GetUserIDFromContext(ctx *gin.Context) (uint, bool) {
	userID := ctx.GetUint("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "user id not found in context",
		})
		return 0, false
	}
	return userID, true
}
