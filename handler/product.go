package handler

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/gorm"
)

type ProductData struct {
	Code   string   `json:"code"`
	Name   string   `json:"name"`
	Price  float64  `json:"price"`
	Weight *float64 `json:"weight"`
}

func GetProducts() (*[]ProductData, error) {
	db := database.DB
	var products []models.Product

	if err := db.Select("code", "name", "price", "weight").Find(&products).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	productData := make([]ProductData, len(products))
	for i, p := range products {
		productData[i] = ProductData{
			Code:   p.Code,
			Name:   p.Name,
			Price:  p.Price,
			Weight: p.Weight,
		}
	}

	return &productData, nil
}

func GetProductByCode(code string) (*ProductData, error) {
	db := database.DB
	var product models.Product
	if err := db.Where(&models.Product{Code: strings.ToUpper(code)}).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	productData := ProductData{
		Code:   product.Code,
		Name:   product.Name,
		Price:  product.Price,
		Weight: product.Weight,
	}

	return &productData, nil
}

func GetAllProducts(c *fiber.Ctx) error {
	products, err := GetProducts()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get products", "data": err})
	}

	if products == nil || len(*products) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No products found", "data": nil})
	}

	return c.Status(200).JSON(products)
}

func GetProduct(c *fiber.Ctx) error {
	code := c.Params("code")

	product, err := GetProductByCode(code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get product", "data": err})
	}

	if product == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid product code", "data": nil})
	}

	return c.Status(200).JSON(product)
}
