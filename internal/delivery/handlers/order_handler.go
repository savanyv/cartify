package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type OrderHandler struct {
	orderUsecase *usecase.OrderUsecase
	validator    *helpers.ValidatorService
}

func NewOrderHandler(ou *usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: ou,
		validator:    helpers.NewValidatorService(),
	}
}

// ==================== USER METHODS ====================

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	order, err := h.orderUsecase.CreateOrder(c.Context(), userID)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.SuccessCreated(c, "Order created successfully", order)
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := h.orderUsecase.GetuserOrders(c.Context(), userID, page, limit, search, sort, order)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.SuccessPaginated(c, "Orders retrieved successfully", items, total, page, limit)
}

func (h *OrderHandler) GetOrderByID(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	orderID := c.Params("id")

	order, err := h.orderUsecase.GetOrderByID(c.Context(), userID, orderID)
	if err != nil {
		return helpers.NotFound(c, err.Error())
	}

	return helpers.Success(c, "Order retrieved successfully", order)
}

// ==================== ADMIN METHODS ====================

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	items, total, err := h.orderUsecase.GetAllOrders(c.Context(), page, limit, search, sort, order)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.SuccessPaginated(c, "All orders retrieved successfully", items, total, page, limit)
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var req dto.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	if err := h.orderUsecase.UpdateOrderStatus(c.Context(), orderID, model.OrderStatus(req.Status)); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Order status updated successfully", nil)
}