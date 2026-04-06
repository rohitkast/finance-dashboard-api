package api

import (
	"errors"
	"log"
	"net/http"
	"personal_finance_dashboard/config"
	"personal_finance_dashboard/internal/models"
	"personal_finance_dashboard/internal/repository"
	"personal_finance_dashboard/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CreateUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin viewer analyst"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userInput CreateUserInput

		if err := ctx.ShouldBindJSON(&userInput); err != nil {
			log.Printf("CreateUser: invalid request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		if len(userInput.Password) < 8 {
			log.Printf("CreateUser: password too short for email=%s", userInput.Email)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Password should be atleast 8 characters"})
			return
		}

		// password hashing
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(userInput.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("CreateUser: failed to hash password for email=%s: %v", userInput.Email, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		user := &models.User{
			Email:    userInput.Email,
			Password: string(hashedPass),
			IsActive: true,
			Role:     userInput.Role,
		}

		err = repository.CreateUser(db, user)
		if err != nil {
			log.Printf("CreateUser: failed to create user email=%s: %v", userInput.Email, err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User created successfully",
		})

	}
}

func LoginUser(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginInput LoginUserInput

		if err := ctx.ShouldBindJSON(&loginInput); err != nil {
			log.Printf("LoginUser: invalid request body: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}

		// fetch user by email
		user, err := repository.GetUserByEmail(db, loginInput.Email)
		if err != nil {
			log.Printf("LoginUser: invalid credentials. No user with this email exist: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "No user with this email or invalid credentials"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password))
		if err != nil {
			log.Printf("LoginUser: invalid credentials. Password incorrect: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid credentials. Password incorrect"})
			return
		}

		// claims are payload to jwtNewWithClaims
		claims := jwt.MapClaims{
			"user_id":   user.ID,
			"email":     user.Email,
			"user_role": user.Role,
			"exp":       time.Now().Add(24 * time.Hour).Unix(), //expiry
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			log.Printf("LoginUser: jwt token error %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "failed to generate token " + err.Error()})
			return
		}

		// successfully loguser
		ctx.JSON(http.StatusOK, LoginResponse{AccessToken: tokenString})
	}
}

func GetUserById(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("GetUserById: invalid id=%q", idStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		user, err := repository.GetUserById(db, uint(id))
		if err != nil {
			log.Printf("GetUserById: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User fetched successfully",
			"data":    user,
		})
	}
}

func DeleteUserById(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, ok := utils.GetUserIDFromContext(ctx)
		if !ok {
			return
		}

		idStr := ctx.Param("id")
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil || id == 0 {
			log.Printf("DeleteUserById: invalid id=%q", idStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid user id"})
			return
		}

		err = repository.DeleteUser(db, uint(id))
		if err != nil {
			log.Printf("DeleteUserById: failed for id=%d: %v", id, err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User deleted successfully",
		})
	}
}
