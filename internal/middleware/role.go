package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorizedRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exists := ctx.Get("user_role")
		if !exists {
			log.Printf("RoleMiddleware: missing user role claims=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User_role not found in context"})
			ctx.Abort()
			return
		}

		userRole := value.(string)
		authorized := false

		// for mentioned allowed roles is he in that role?
		for _, role := range allowedRoles {
			if role == userRole {
				authorized = true
				break
			}
		}

		if !authorized {
			log.Printf("RoleMiddleware: user not authorized to use these routes=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "this user is not allowed here"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
