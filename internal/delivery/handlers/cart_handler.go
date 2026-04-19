package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type CartHandler struct {
	cartUsecase *usecase.CartUsecase
	validator   *helpers.ValidatorService
}

func NewCartHandler(cu *usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cu,
		validator:   helpers.NewValidatorService(),
	}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	cart, err := h.cartUsecase.GetCart(c.Context(), userID)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Cart retrieved successfully", cart)
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req dto.AddToCartRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	if err := h.cartUsecase.AddToCart(c.Context(), userID, req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Item added to cart successfully", nil)
}

func (h *CartHandler) UpdateCartItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	itemID := c.Params("itemId")

	var req dto.UpdateCartItemRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	if err := h.cartUsecase.UpdateCartItem(c.Context(), userID, itemID, req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Cart item updated successfully", nil)
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	itemID := c.Params("itemId")

	if err := h.cartUsecase.RemoveFromCart(c.Context(), userID, itemID); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Item removed from cart successfully", nil)
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.cartUsecase.ClearCart(c.Context(), userID); err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Cart cleared successfully", nil)
}