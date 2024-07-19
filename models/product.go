package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code     string   `json:"code" gorm:"unique;not null"`
	Name     string   `json:"name" gorm:"not null"`
	Price    float64  `json:"price" gorm:"type:numeric(10,2);not null"`
	Weight   *float64 `json:"weight" gorm:"type:numeric(3,2)"`
	Merchant string   `json:"merchant" gorm:"not null"`

	User User `gorm:"foreignKey:Merchant;references:Username"`
}

type CreateProductValidation struct {
	Code   string   `json:"code" validate:"required,min=3"`
	Name   string   `json:"name" validate:"required,min=3"`
	Price  float64  `json:"price" validate:"required,gt=0"`
	Weight *float64 `json:"weight" validate:"gt=0"`
}

type UpdateProductValidation struct {
	Name   string   `json:"name" validate:"required,min=3"`
	Price  float64  `json:"price" validate:"required,gt=0"`
	Weight *float64 `json:"weight" validate:"gt=0"`
}
