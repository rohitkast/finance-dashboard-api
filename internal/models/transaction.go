package models

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	UserID      uint    `gorm:"not null"`
	Amount      float64 `gorm:"not null"`
	Description string
	Category    string `gorm:"check:category IN ('income','expense')" json:"category"`
	User        User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}
