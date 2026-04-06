package routes

import (
	"personal_finance_dashboard/api"
	"personal_finance_dashboard/config"
	"personal_finance_dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TransactionsRoutes(r *gin.RouterGroup, db *gorm.DB) {
	router := r.Group("/transaction")
	router.Use(middleware.AuthMiddleware(&config.Config{}))

	// Read operations - all authenticated users
	readRoutes := router.Group("")
	readRoutes.Use(middleware.AuthorizedRoles("admin", "viewer", "analyst"))
	{
		readRoutes.GET("/all", api.GetAllTransactions(db))
		readRoutes.GET("/:id", api.GetTransactionById(db))
		readRoutes.GET("/filter", api.GetFilteredTransactions(db))
	}

	// Write operations - only analyst and admin
	writeRoutes := router.Group("")
	writeRoutes.Use(middleware.AuthorizedRoles("admin", "analyst"))
	{
		writeRoutes.POST("/create", api.CreateTransaction(db))
		writeRoutes.PUT("/:id", api.UpdateTransaction(db))
		writeRoutes.DELETE("/:id", api.DeleteTransaction(db))
	}
}
