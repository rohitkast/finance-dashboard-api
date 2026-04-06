package middleware

import (
	"fmt"
	"log"
	"net/http"
	"personal_finance_dashboard/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("AuthMiddleware: missing Authorization header path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" || tokenString == authHeader {
			log.Printf("AuthMiddleware: invalid bearer token format path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing bearer token"})
			ctx.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// checking if the algorithum used is same. if yes, we give it the jwtsecret to check further
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpection signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil {
			log.Printf("AuthMiddleware: token parse failed path=%s ip=%s: %v", ctx.FullPath(), ctx.ClientIP(), err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		// now to extract claims (user id, email)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("AuthMiddleware: invalid token claims path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			log.Printf("AuthMiddleware: missing or invalid user_id claim path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		userRole, ok := claims["user_role"].(string)
		if !ok {
			log.Printf("AuthMiddleware: missing or invalid user_role claim path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			ctx.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			expiration_time := time.Unix(int64(exp), 0)

			// so if time now is after expiration. its expred
			if time.Now().After(expiration_time) {
				log.Printf("AuthMiddleware: token expired path=%s ip=%s", ctx.FullPath(), ctx.ClientIP())
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				ctx.Abort()
				return
			}
		}

		// everything good. set typed values on the gin context
		ctx.Set("user_id", uint(userID))
		ctx.Set("user_role", userRole)
		log.Printf("user id and user role set as %v and %v:", userID, userRole)
		// move forward with the handler or next middleware
		ctx.Next()
	}
}
