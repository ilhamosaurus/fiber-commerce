package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code   string  `json:"code" gorm:"uniqueIndex;not null"`
	Name   string  `json:"name" gorm:"not null"`
	Price  float64 `json:"price" gorm:"type:numeric(10,2);not null"`
	Weight float64 `json:"weight" gorm:"type:numeric(3,2)"`
}
