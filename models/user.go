package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string   `json:"username" gorm:"unique;not null"`
	Password string   `json:"password" gorm:"not null"`
	Account  *Account `gorm:"foreignKey:Owner;references:Username"`
}

type UserValidation struct {
	Username string `json:"username" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
}
