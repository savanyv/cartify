package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type CartHandler struct {
	cartUsecase *usecase.CartUsecase
	validator *helpers.ValidatorService
}

func NewCartHandler(cu *usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cu,
		validator: helpers.NewValidatorService(),
	}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	cart, err := h.cartUsecase.GetCart(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart retrieved successfully",
		"data":    cart,
	})
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req dto.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.cartUsecase.AddToCart(c.Context(), userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product added to cart successfully",
	})
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	itemID := c.Params("item_id")

	var req dto.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.cartUsecase.UpdateCartItem(c.Context(), userID, itemID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart item updated successfully",
	})
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	itemID := c.Params("item_id")

	if err := h.cartUsecase.RemoveFromCart(c.Context(), userID, itemID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product removed from cart successfully",
	})
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.cartUsecase.ClearCart(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Cart cleared successfully",
	})
}