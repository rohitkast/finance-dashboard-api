package routes

import (
	"personal_finance_dashboard/api"
	"personal_finance_dashboard/config"
	"personal_finance_dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRoutes(r *gin.RouterGroup, db *gorm.DB) {
	router := r.Group("/admin")
	router.Use(middleware.AuthMiddleware(&config.Config{}))
	router.Use(middleware.AuthorizedRoles("admin"))
	{
		router.GET("/users", api.GetAllUsers(db))
		router.GET("/getuser/:id", api.GetUserById(db))
		router.DELETE("/user/:id", api.DeleteUserById(db))
		router.GET("/transactions/:uid", api.AdminGetUserTransactions(db))
		router.DELETE("/transaction/:id", api.AdminDeleteTransaction(db))
	}
}
