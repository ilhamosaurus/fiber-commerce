package handler

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ilhamosaurus/fiber-commerce/database"
	"github.com/ilhamosaurus/fiber-commerce/models"
	"gorm.io/gorm"
)

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
