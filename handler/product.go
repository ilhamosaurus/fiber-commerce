package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"github.com/ilhamosaurus/fiber-commerce/util"
	"gorm.io/gorm"
)

type ProductData struct {
	Code     string   `json:"code"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Weight   *float64 `json:"weight"`
	Merchant string   `json:"merchant"`
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
			Code:     p.Code,
			Name:     p.Name,
			Price:    p.Price,
			Weight:   p.Weight,
			Merchant: p.Merchant,
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
		Code:     product.Code,
		Name:     product.Name,
		Price:    product.Price,
		Weight:   product.Weight,
		Merchant: product.Merchant,
	}

	return &productData, nil
}

// @Summary Get all products
// @Description Get all products
// @Tags Products
// @Produce json
// @Success 200 {array} handler.ProductData	"OK"
// @Failure 404 {object} string "No products found"
// @Failure 500 {object} string "Failed to get products"
// @Router /api/products [get]
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

// @Summary Get product by code
// @Description Get product by code
// @Tags Products
// @Produce json
// @Param code path string true "Product code"
// @Success 200 {object} handler.ProductData	"OK"
// @Failure 404 {object} string "Invalid product code"
// @Failure 500 {object} string "Failed to get product"
// @Router /api/products/{code} [get]
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

// @Summary Create product
// @Description Create product
// @Tags Products
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body models.CreateProductValidation true "Product data"
// @Success 201 {object} handler.ProductData "Product created successfully"
// @Failure 400 {object} string "Invalid fields"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Failed to create product"
// @Router /api/products [post]
func CreateProduct(c *fiber.Ctx) error {
	user := util.CurrentUser(c)
	if user.Role != models.Merchant {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	db := database.DB
	var body models.CreateProductValidation
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields"})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
	}

	if err := db.Create(&models.Product{Code: strings.ToUpper(body.Code), Name: body.Name, Price: body.Price, Weight: body.Weight, Merchant: user.Username}).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.Status(409).JSON(fiber.Map{"error": "Product's code already exists"})
		}
	}

	response := ProductData{
		Code:     body.Code,
		Name:     body.Name,
		Price:    body.Price,
		Weight:   body.Weight,
		Merchant: user.Username,
	}

	return c.Status(201).JSON(response)
}

// @Summary Update product
// @Description Update product
// @Tags Products
// @Security Bearer
// @Accept json
// @Produce json
// @Param code path string true "Product code"
// @Param body body models.UpdateProductValidation true "Product data"
// @Success 200 {object} handler.ProductData "Product updated successfully"
// @Failure 400 {object} string "Invalid fields"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Invalid product code"
// @Failure 500 {object} string "Failed to update product"
// @Router /api/products/{code} [put]
func UpdateProduct(c *fiber.Ctx) error {
	user := util.CurrentUser(c)
	if user.Role != models.Merchant {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	code := c.Params("code")
	db := database.DB

	product, err := GetProductByCode(code)
	if err != nil || product == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid Product Code", "data": err})
	}

	if product.Merchant != user.Username {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var body models.UpdateProductValidation
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid fields", "data": err})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
	}

	if err := db.Model(&models.Product{}).Where("code = ?", product.Code).Updates(models.Product{Name: body.Name, Price: body.Price, Weight: body.Weight}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update product", "data": err})
	}

	response := ProductData{
		Code:     product.Code,
		Name:     body.Name,
		Price:    body.Price,
		Weight:   body.Weight,
		Merchant: user.Username,
	}

	return c.Status(200).JSON(response)
}

// @Summary Delete product
// @Description Delete product
// @Tags Products
// @Security Bearer
// @Accept json
// @Produce json
// @Param code path string true "Product code"
// @Success 200 {object} string "Product deleted successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Invalid product code"
// @Failure 500 {object} string "Failed to delete product"
// @Router /api/products/{code} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	user := util.CurrentUser(c)
	if user.Role != models.Merchant {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}
	code := c.Params("code")
	db := database.DB

	product, err := GetProductByCode(code)
	if err != nil || product == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Invalid Product Code", "data": err})
	}

	if product.Merchant != user.Username {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := db.Delete(&product).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete product", "data": err})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Product deleted successfully"})
}
