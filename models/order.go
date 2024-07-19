package models

import (
	"database/sql/driver"
	"fmt"

	"gorm.io/gorm"
)

type Type string

const (
	Topup   Type = "TOPUP"
	Payment Type = "PAYMENT"
	Revenue Type = "REVENUE"
)

func (t *Type) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*t = Type(v)
	case string:
		*t = Type(v)
	default:
		return fmt.Errorf("expected []byte or string, got %T", value)
	}
	return nil
}

func (t Type) Value() (driver.Value, error) {
	return string(t), nil
}

type Order struct {
	gorm.Model
	Invoice     string  `json:"invoice" gorm:"unique;not null"`
	AccountID   uint    `json:"account_id" gorm:"not null"`
	Merchant    *string `json:"merchant" `
	Buyer       *string `json:"buyer" `
	Amount      float64 `json:"amount" gorm:"type:numeric(10,2);not null"`
	Type        Type    `json:"type" gorm:"not null; type:order_type"`
	Description *string `json:"description" gorm:"type:text"`

	Account Account `gorm:"foreignKey:AccountID;references:ID"`
}

type TopupValidation struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

type PaymentValidation struct {
	Code string `json:"code" validate:"required"`
	Qty  int    `json:"qty" validate:"required,gt=0"`
}
