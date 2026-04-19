package helpers

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Error      interface{} `json:"error,omitempty"`
}

type PaginationMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	HasPrev    bool  `json:"has_prev"`
	HasNext    bool  `json:"has_next"`
}

func SuccessPaginated(c *fiber.Ctx, message string, data interface{}, total int64, page, limit int) error {
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasPrev := page > 1
	hasNext := page < totalPages

	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: PaginationMeta{
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			HasPrev:    hasPrev,
			HasNext:    hasNext,
		},
	})
}

func Success(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func SuccessCreated(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Shortcut functions
func BadRequest(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusBadRequest, message, nil)
}

func Unauthorized(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusUnauthorized, message, nil)
}

func Forbidden(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusForbidden, message, nil)
}

func NotFound(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusNotFound, message, nil)
}

func Conflict(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusConflict, message, nil)
}

func InternalServerError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, fiber.StatusInternalServerError, message, nil)
}

func ValidationError(c *fiber.Ctx, errors interface{}) error {
	return ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", errors)
}