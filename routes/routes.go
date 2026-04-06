package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(app *gin.Engine, db *gorm.DB) {
	api := app.Group("/api")

	TransactionsRoutes(api, db)
	DashboardRoutes(api, db)
	UserRoutes(api, db)
	AdminRoutes(api, db)
}
