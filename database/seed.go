package database

import (
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/gorm"
)

var products = []models.Product {
	{
    Code: "PAJAK",
    Name: "Pajak PBB",
    Price: 40000,
  },
  {
    Code: "PLN",
    Name: "Listrik",
    Price: 10000,
  },
  {
    Code: "PDAM",
    Name: "PDAM Berlangganan",
    Price: 40000,
  },
  {
    Code: "PULSA",
    Name: "Pulsa",
    Price: 40000,
  },
  {
    Code: "PGN",
    Name: "PGN Berlangganan",
    Price: 50000,
  },
  {
    Code: "MUSIK",
    Name: "Musik Berlangganan",
    Price: 50000,
  },
  {
    Code: "TV",
    Name: "TV Berlangganan",
    Price: 50000,
  },
  {
    Code: "PAKET_DATA",
    Name: "Paket data",
    Price: 50000,
  },
  {
    Code: "VOUCHER_GAME",
    Name: "Voucher Game",
    Price: 100000,
  },
  {
    Code: "VOUCHER_MAKANAN",
    Name: "Voucher Makanan",
    Price: 100000,
  },
  {
    Code: "ZAKAT",
    Name: "Zakat",
    Price: 300000,
  },
}

func Load(db *gorm.DB) {
	db.Create(&products)
}