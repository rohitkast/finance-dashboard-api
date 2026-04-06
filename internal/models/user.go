package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Role     string `gorm:"check:role IN ('admin','viewer','analyst')" json:"role"`
}
