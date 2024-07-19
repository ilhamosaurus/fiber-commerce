package database

import (
	"log"

	"github.com/ilhamosaurus/fiber-commerce/models"
	"github.com/ilhamosaurus/fiber-commerce/util"
	"gorm.io/gorm"
)

var products = []models.Product{
	{
		Code:     "PAJAK",
		Name:     "Pajak PBB",
		Price:    40000,
		Merchant: "admin",
	},
	{
		Code:     "PLN",
		Name:     "Listrik",
		Price:    10000,
		Merchant: "admin",
	},
	{
		Code:     "PDAM",
		Name:     "PDAM Berlangganan",
		Price:    40000,
		Merchant: "admin",
	},
	{
		Code:     "PULSA",
		Name:     "Pulsa",
		Price:    40000,
		Merchant: "admin",
	},
	{
		Code:     "PGN",
		Name:     "PGN Berlangganan",
		Price:    50000,
		Merchant: "admin",
	},
	{
		Code:     "MUSIK",
		Name:     "Musik Berlangganan",
		Price:    50000,
		Merchant: "admin",
	},
	{
		Code:     "TV",
		Name:     "TV Berlangganan",
		Price:    50000,
		Merchant: "admin",
	},
	{
		Code:     "PAKET_DATA",
		Name:     "Paket data",
		Price:    50000,
		Merchant: "admin",
	},
	{
		Code:     "VOUCHER_GAME",
		Name:     "Voucher Game",
		Price:    100000,
		Merchant: "admin",
	},
	{
		Code:     "VOUCHER_MAKANAN",
		Name:     "Voucher Makanan",
		Price:    100000,
		Merchant: "admin",
	},
	{
		Code:     "ZAKAT",
		Name:     "Zakat",
		Price:    300000,
		Merchant: "admin",
	},
}

func Load(db *gorm.DB) {
	hash, err := util.HashedPassword("qwerty")
	if err != nil {
		log.Fatal(err)
	}
	user := models.User{
		Username: "admin",
		Password: hash,
		Role:     models.Role("MERCHANT"),
		Account: &models.Account{
			Owner:   "admin",
			Balance: 0,
		},
	}

	db.Create(&user)
	db.Create(&products)
}
