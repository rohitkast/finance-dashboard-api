package routes

import (
	"personal_finance_dashboard/api"
	"personal_finance_dashboard/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserRoutes(r *gin.RouterGroup, db *gorm.DB) {
	authRouter := r.Group("/auth")
	{
		authRouter.POST("/create", api.CreateUser(db))
		authRouter.POST("/login", api.LoginUser(db, &config.Config{}))
	}
}
