package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/savanyv/cartify/internal/utils/helpers"
)

type ProductHandler struct {
	productUsecase *usecase.ProductUsecase
	validator      *helpers.ValidatorService
}

func NewProductHandler(pu *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: pu,
		validator:      helpers.NewValidatorService(),
	}
}

// ================== Product User ==================

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	search := c.Query("search", "")
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	// Validate
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Usecase returns: (items, total, error)
	items, total, err := h.productUsecase.GetAllProductsWithPagination(c.Context(), page, limit, search, sort, order)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.SuccessPaginated(c, "Products retrieved successfully", items, total, page, limit)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.productUsecase.GetProductByID(c.Context(), id)
	if err != nil {
		return helpers.NotFound(c, err.Error())
	}

	return helpers.Success(c, "Product retrieved successfully", product)
}

// ================== Product Admin ==================

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	product, err := h.productUsecase.CreateProduct(c.Context(), req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.SuccessCreated(c, "Product created successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	product, err := h.productUsecase.UpdateProduct(c.Context(), id, req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Product updated successfully", product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.productUsecase.DeleteProduct(c.Context(), id); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Product deleted successfully", nil)
}

// ================== Variant Admin ==================

func (h *ProductHandler) CreateVariant(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	var req dto.CreateVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	variant, err := h.productUsecase.CreateVariant(c.Context(), productID, req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.SuccessCreated(c, "Variant created successfully", variant)
}

func (h *ProductHandler) UpdateVariant(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return helpers.ValidationError(c, err.Error())
	}

	variant, err := h.productUsecase.UpdateVariant(c.Context(), id, req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Variant updated successfully", variant)
}