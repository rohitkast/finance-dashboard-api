package routes

import (
	"personal_finance_dashboard/api"
	"personal_finance_dashboard/config"
	"personal_finance_dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DashboardRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// dashboard routes are where response have calculated fields like total, balance etc
	router := r.Group("/dashboard")
	router.Use(middleware.AuthMiddleware(&config.Config{}))
	router.Use(middleware.AuthorizedRoles("admin", "viewer", "analyst"))
	{
		router.GET("/summary", api.GetSummary(db))
		router.GET("/recent", api.GetRecent(db))
		router.GET("/category/:category", api.GetTransactionsByCategory(db))
		router.GET("/trends/:category", api.GetMonthlyTrends(db))
	}
}
