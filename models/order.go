package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type Type string

const (
	Topup   Type = "TOPUP"
	Payment Type = "PAYMENT"
	Revenue Type = "REVENUE"
)

func (t *Type) Scan(value interface{}) error {
	*t = Type(value.([]byte))
	return nil
}

func (t Type) Value() (driver.Value, error) {
	return string(t), nil
}

type Order struct {
	gorm.Model
	Invoice   string  `json:"invoice" gorm:"uniqueIndex;not null"`
	AccountID uint    `json:"account_id" gorm:"not null"`
	Merchant  string  `json:"merchant" `
	Buyer     string  `json:"buyer" `
	Amount    float64 `json:"amount" gorm:"type:numeric(10,2);not null"`
	Type      Type    `json:"type" gorm:"not null; type:order_type"`

	Account Account `gorm:"foreignKey:AccountID;references:ID"`
}
