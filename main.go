package main

import (
	"log"
	"personal_finance_dashboard/config"
	"personal_finance_dashboard/database"
	"personal_finance_dashboard/internal/models"
	"personal_finance_dashboard/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	var db *gorm.DB
	var err error
	db, err = database.Connect(cfg.DatabaseURL)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// geting underlying SQL DB for proper cleanup
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB:", err)
	}
	defer sqlDB.Close()

	// auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// setup routes
	app := gin.Default()
	routes.RegisterRoutes(app, db)

	log.Printf("Starting server on port %s", cfg.Port)
	if err := app.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
