package models

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	Owner   string  `json:"owner" gorm:"uniqueIndex;not null"`
	Balance float64 `json:"balance" gorm:"type:numeric(10,2);not null"`
	User    User    `gorm:"foreignKey:Owner;references:Username"`
}
