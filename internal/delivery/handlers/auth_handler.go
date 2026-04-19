package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type AuthHandler struct {
	authUsecase *usecase.AuthUsecase
	validator   *helpers.ValidatorService
}

func NewAuthHandler(au *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		authUsecase: au,
		validator:   helpers.NewValidatorService(),
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	resp, err := h.authUsecase.Register(c.Context(), req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.SuccessCreated(c, "Registration successful", resp)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	resp, err := h.authUsecase.Login(c.Context(), req)
	if err != nil {
		return helpers.Unauthorized(c, err.Error())
	}

	return helpers.Success(c, "Login successful", resp)
}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	user, err := h.authUsecase.GetUserByID(c.Context(), userID)
	if err != nil {
		return helpers.NotFound(c, err.Error())
	}

	return helpers.Success(c, "Profile retrieved", user)
}

func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	if err := h.authUsecase.ChangePassword(c.Context(), userID, req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Password changed successfully", nil)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	if err := h.authUsecase.Logout(c.Context(), userID); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Logout successful", nil)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	resp, err := h.authUsecase.RefreshToken(c.Context(), req.RefreshToken)
	if err != nil {
		return helpers.Unauthorized(c, err.Error())
	}

	return helpers.Success(c, "Token refreshed", resp)
}
