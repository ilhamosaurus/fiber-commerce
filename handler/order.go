package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/gorm"
)

// @Summary Topup user's balance
// @Tags Transaction
// @Description Topup user's balance
// @Security Bearer
// @Accept json
// @Produce json
// @Param body body models.TopupValidation true "Topup"
// @Success 201 {object} handler.Topup.TopupResponse "OK"
// @Failure 400 {object} string "Invalid fields"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Failed to topup"
// @Router /api/transaction/topup [post]
func Topup(c *fiber.Ctx) error {
	type TopupResponse struct {
		Invoice   string      `json:"invoice"`
		Amount    float64     `json:"amount"`
		Type      models.Type `json:"type"`
		CreatedAt time.Time   `json:"created_at"`
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	db := database.DB
	var err error

	body := &models.TopupValidation{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Fields"})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
	}

	transactionUtil, err := GetInvNumber(username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to topup", "data": err})
	}
	if transactionUtil == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
	}
	transaction := models.Order{
		AccountID: transactionUtil.Account.ID,
		Invoice:   transactionUtil.Invoice,
		Amount:    body.Amount,
		Type:      models.Type("TOPUP"),
	}

	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	tx.Model(&models.Account{}).Where("owner = ?", username).Update("balance", gorm.Expr("balance + ?", body.Amount))
	if err := tx.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to topup", "data": err})
	}

	response := TopupResponse{
		Invoice:   transactionUtil.Invoice,
		Amount:    body.Amount,
		Type:      models.Type("TOPUP"),
		CreatedAt: transaction.CreatedAt,
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "success topup",
		"data":    response,
	})
}

// @Summary Get user's transactions
// @Tags Transaction
// @Description Get user's transactions history
// @Security Bearer
// @Produce json
// @Success 200 {object} handler.GetOrders.OrderResponse
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "No transactions found"
// @Failure 500 {object} string "Failed to get transactions"
// @Router /api/transaction/history [get]
func GetOrders(c *fiber.Ctx) error {
	type OrderResponse struct {
		Invoice   string      `json:"invoice"`
		Merchant  *string     `json:"merchant"`
		Buyer     *string     `json:"buyer"`
		Amount    float64     `json:"amount"`
		Type      models.Type `json:"type"`
		CreatedAt time.Time   `json:"created_at"`
	}

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	db := database.DB

	account, err := GetAccountByUsername(username)
	if err != nil || account == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Failed to get account", "data": err})
	}

	orders := []models.Order{}
	if err := db.Where("account_id = ?", account.ID).Order("created_at DESC").Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "No transactions found", "data": nil})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get transactions", "data": err})
	}

	orderResponse := make([]OrderResponse, len(orders))
	for i, order := range orders {
		orderResponse[i] = OrderResponse{
			Invoice:   order.Invoice,
			Merchant:  order.Merchant,
			Buyer:     order.Buyer,
			Amount:    order.Amount,
			Type:      order.Type,
			CreatedAt: order.CreatedAt,
		}
	}

	return c.Status(200).JSON(fiber.Map{"data": orderResponse})
}

// @Summary Payment
// @Tags Transaction
// @Description User's purchase products and pay the merchant
// @Security Bearer
// @Accept json
// @Produce json
// @Param payment body models.PaymentValidation true "Payment"
// @Success 201 {object} handler.Payment.PaymentResponse
// @Failure 400 {object} string "Invalid Fields"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Product not found"
// @Failure 500 {object} string "Failed to payment"
// @Router /api/transaction/payment [post]
func Payment(c *fiber.Ctx) error {
	type PaymentResponse struct {
		Invoice   string      `json:"invoice"`
		Merchant  *string     `json:"merchant"`
		Buyer     *string     `json:"buyer"`
		Amount    float64     `json:"amount"`
		Type      models.Type `json:"type"`
		CreatedAt time.Time   `json:"created_at"`
	}
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	db := database.DB

	body := &models.PaymentValidation{}
	if err := c.BodyParser(body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Fields"})
	}

	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		errMsgs := make([]string, 0)
		for _, e := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s %s", e.Field(), e.Tag(), e.Param()))
		}
		return c.Status(400).JSON(fiber.Map{"error": errMsgs})
	}

	transactionUtil, err := GetInvNumber(username)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to purchase", "data": err})
	}
	if transactionUtil == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Account not found"})
	}

	product, err := GetProductByCode(body.Code)
	if err != nil || product == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Product not found", "data": err})
	}
	// balance check
	totalAmount := product.Price * float64(body.Qty)
	if totalAmount > transactionUtil.Account.Balance {
		return c.Status(400).JSON(fiber.Map{"error": "Insufficient balance"})
	}

	merchatUtil, err := GetInvNumber(product.Merchant)
	if err != nil || merchatUtil == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Merchant not found", "data": err})
	}

	description := fmt.Sprintf("Payment for product %s(%s)", product.Name, product.Code)
	merchantTransaction := models.Order{
		AccountID:   merchatUtil.Account.ID,
		Invoice:     merchatUtil.Invoice,
		Amount:      product.Price,
		Type:        models.Type("REVENUE"),
		Merchant:    &product.Merchant,
		Buyer:       &username,
		Description: &description,
	}

	transaction := models.Order{
		AccountID:   transactionUtil.Account.ID,
		Invoice:     transactionUtil.Invoice,
		Amount:      totalAmount,
		Type:        models.Type("PAYMENT"),
		Merchant:    &product.Merchant,
		Buyer:       &username,
		Description: &description,
	}

	tx := db.Session(&gorm.Session{SkipDefaultTransaction: true})
	if err := tx.Model(&models.Account{}).Where("owner = ?", username).Update("balance", gorm.Expr("balance - ?", totalAmount)).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to purchase", "data": err})
	}
	if err := tx.Model(&models.Account{}).Where("owner = ?", product.Merchant).Update("balance", gorm.Expr("balance + ?", totalAmount)).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to purchase", "data": err})
	}
	if err := tx.Create(&transaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to purchase", "data": err})
	}
	if err := tx.Create(&merchantTransaction).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to purchase", "data": err})
	}

	paymentResponse := PaymentResponse{
		Invoice:   transaction.Invoice,
		Merchant:  transaction.Merchant,
		Buyer:     transaction.Buyer,
		Amount:    transaction.Amount,
		Type:      transaction.Type,
		CreatedAt: transaction.CreatedAt,
	}

	return c.Status(201).JSON(fiber.Map{"data": paymentResponse})
}
