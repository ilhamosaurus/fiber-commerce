package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ilhamosaurus/fiber-commerce/config"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDb() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("failed to parse database port")
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// init enum for role and order
	db.Exec("CREATE TYPE role AS ENUM ('CLIENT', 'MERCHANT')")
	db.Exec("CREATE TYPE order_type AS ENUM ('TOPUP', 'PAYMENT', 'REVENUE')")

	fmt.Println("Connection Opened to Database")
	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Account{}, &models.Order{})
	fmt.Println("Database Migrated")
	Load(db) // products seeding
	DB = db
}
