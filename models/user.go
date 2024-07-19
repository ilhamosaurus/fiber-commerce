package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

type Role string

const (
	Client   Role = "CLIENT"
	Merchant Role = "MERCHANT"
)

func (t *Role) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*t = Role(v)
	case string:
		*t = Role(v)
	default:
		return fmt.Errorf("expected []byte or string, got %T", value)
	}
	return nil
}

func (t Role) Value() (driver.Value, error) {
	return string(t), nil
}

type User struct {
	gorm.Model
	Username string   `json:"username" gorm:"unique;not null"`
	Password string   `json:"password" gorm:"not null"`
	Role     Role     `json:"role" gorm:"not null; type:role"`
	
	Account  *Account `gorm:"foreignKey:Owner;references:Username"`
}

type RegisterValidation struct {
	Username string `json:"username" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
	Role     Role   `json:"role" validate:"required,role"`
}

type LoginValidation struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
