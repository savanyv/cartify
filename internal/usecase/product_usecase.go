package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
)

type ProductUsecase struct {
	productRepo model.ProductRepository
	productVariantRepo model.ProductVariantRepository
}

func NewProductUsecase(pr model.ProductRepository, pvr model.ProductVariantRepository) *ProductUsecase {
	return &ProductUsecase{
		productRepo: pr,
		productVariantRepo: pvr,
	}
}

// ========================== Product Methods (Public) ==========================
func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := u.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductResponse
	for _, product := range products {
		responses = append(responses, u.toProductResponse(&product))
	}

	return responses, nil
}

func (u *ProductUsecase) GetProductByID(ctx context.Context, ID string) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	resp := u.toProductResponse(product)
	return &resp, nil
}

// ========================== Product Methods (Admin) ==========================
func (u *ProductUsecase) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	product := &model.Product{
		Name: req.Name,
		Description: req.Description,
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	resp := u.toProductResponse(product)
	return &resp, nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, ID string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	product, err := u.productRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}

	if err := u.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	product, _ = u.productRepo.FindByID(ctx, ID)
	resp := u.toProductResponse(product)

	return &resp, nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, ID string) error {
	product, err := u.productRepo.FindByID(ctx, ID)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return u.productRepo.Delete(ctx, ID)
}

// ========================== Product Variant Methods (Admin) ==========================
func (u *ProductUsecase) CreateVariant(ctx context.Context, productID string, req dto.CreateVariantRequest) (*dto.VariantResponse, error) {
	product, err := u.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	parsedProductID, err := uuid.Parse(productID)
	if err != nil {
		return nil, err
	}

	variant := &model.ProductVariant{
		ProductID: parsedProductID,
		Name: req.Name,
		Stock: req.Stock,
		Price: req.Price,
	}

	if err := u.productVariantRepo.Create(ctx, variant); err != nil {
		return nil, err
	}

	resp := u.toVariantResponse(variant)
	return &resp, nil
}

func (u *ProductUsecase) UpdateVariant(ctx context.Context, ID string, req dto.UpdateVariantRequest) (*dto.VariantResponse, error) {
	variant, err := u.productVariantRepo.FindByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if variant == nil {
		return nil, errors.New("variant not found")
	}

	if req.Name != "" {
		variant.Name = req.Name
	}
	if req.Stock != nil {
		variant.Stock = *req.Stock
	}
	if req.Price > 0 {
		variant.Price = req.Price
	}

	if err := u.productVariantRepo.Update(ctx, variant); err != nil {
		return nil, err
	}

	resp := u.toVariantResponse(variant)
	return &resp, nil
}

// ========================== Helpers Methods ==========================
func (u *ProductUsecase) toProductResponse(p *model.Product) dto.ProductResponse {
	var variants []dto.VariantResponse
	for _, v := range p.Variants {
		variants = append(variants, u.toVariantResponse(&v))
	}

	return dto.ProductResponse{
		ID: p.ID.String(),
		Name: p.Name,
		Description: p.Description,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		Variants: variants,
	}
}

func (u *ProductUsecase) toVariantResponse(v *model.ProductVariant) dto.VariantResponse {
	return dto.VariantResponse{
		ID: v.ID.String(),
		Name: v.Name,
		Stock: v.Stock,
		Price: v.Price,
		ProductID: v.ProductID.String(),
	}
}