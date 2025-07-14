package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/koriebruh/suplyChainTrack/internal/dto"
)

type Handler struct {
}

func SendError(c *fiber.Ctx, statusCode int, err error, message string) error {
	return c.Status(statusCode).JSON(dto.ErrorResponse{
		Error:   err.Error(),
		Message: message,
		Code:    statusCode,
	})
}

func SendSuccess(c *fiber.Ctx, statusCode int, data interface{}, message string) error {
	return c.Status(statusCode).JSON(dto.SuccessResponse{
		Data:    data,
		Message: message,
		Code:    statusCode,
	})
}

type StakeholderHandler interface {
	CreateStakeholder(c *fiber.Ctx) error
	GetStakeholderByEmail(c *fiber.Ctx) error
	GetStakeholder(c *fiber.Ctx) error
	UpdateStakeholder(c *fiber.Ctx) error
	DeleteStakeholder(c *fiber.Ctx) error
	ListStakeholders(c *fiber.Ctx) error
	GetStakeholderStats(c *fiber.Ctx) error
	VerifyStakeholder(c *fiber.Ctx) error
}

type ProductHandler interface {
	CreateProduct(c *fiber.Ctx) error
	GetProduct(c *fiber.Ctx) error
	GetProductBySKU(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	ListProducts(c *fiber.Ctx) error
	GetProductStats(c *fiber.Ctx) error
	GetProductByManufacture(c *fiber.Ctx) error
}

type SupplyChainHandler interface {
	CreateEvent(c *fiber.Ctx) error
	GetEvent(c *fiber.Ctx) error
	UpdateEvent(c *fiber.Ctx) error
	DeleteEvent(c *fiber.Ctx) error
	ListEvents(c *fiber.Ctx) error
	GetProductTrace(c *fiber.Ctx) error
	VerifyEvent(c *fiber.Ctx) error
	GetEventsByProduct(c *fiber.Ctx) error
	GetEventsByStakeholder(c *fiber.Ctx) error
	ValidateEventSequence(c *fiber.Ctx) error
}

type BlockchainHandler interface {
	CreateTransaction(c *fiber.Ctx) error
	GetTransaction(c *fiber.Ctx) error
	GetTransactionByHash(c *fiber.Ctx) error
	UpdateTransaction(c *fiber.Ctx) error
	ListTransactions(c *fiber.Ctx) error
	UpdateTransactionStatus(c *fiber.Ctx) error
	GetTransactionByEvent(c *fiber.Ctx) error
}
