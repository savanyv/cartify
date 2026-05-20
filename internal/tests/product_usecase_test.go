package tests

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupProductUsecase() (*usecase.ProductUsecase, *mocks.ProductRepositoryMock, *mocks.ProductVariantRepositoryMock) {
	pr := new(mocks.ProductRepositoryMock)
	pvr := new(mocks.ProductVariantRepositoryMock)
	return usecase.NewProductUsecase(pr, pvr), pr, pvr
}

func TestProductUsecase_GetAllProductsWithPagination_Success(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	now := time.Now()
	products := []model.Product{
		{ID: uuid.New(), Name: "P1", Description: "D1", CreatedAt: now},
		{ID: uuid.New(), Name: "P2", Description: "D2", CreatedAt: now},
	}

	pr.On("FindAllWithPagination", ctx, 1, 10, "", "created_at", "desc").Return(products, int64(2), nil)

	resp, total, err := uc.GetAllProductsWithPagination(ctx, 1, 10, "", "created_at", "desc")
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, resp, 2)
	assert.Equal(t, "P1", resp[0].Name)
	pr.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_Success(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pid := uuid.New()
	product := &model.Product{ID: pid, Name: "P1", Description: "D1", CreatedAt: time.Now()}

	pr.On("FindByID", ctx, pid.String()).Return(product, nil)

	resp, err := uc.GetProductByID(ctx, pid.String())
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, pid.String(), resp.ID)
	assert.Equal(t, "P1", resp.Name)
	pr.AssertExpectations(t)
}

func TestProductUsecase_GetProductByID_NotFound(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pr.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.GetProductByID(ctx, "x")
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	assert.Nil(t, resp)
	pr.AssertExpectations(t)
}

func TestProductUsecase_CreateProduct_Success(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	req := dto.CreateProductRequest{Name: "New Product", Description: "Desc"}

	pr.On("Create", ctx, mock.MatchedBy(func(p *model.Product) bool {
		return p.Name == req.Name && p.Description == req.Description
	})).Return(nil)

	resp, err := uc.CreateProduct(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Name, resp.Name)
	pr.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_Success(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pid := uuid.New()
	req := dto.UpdateProductRequest{Name: "Updated", Description: "Updated Desc"}

	original := &model.Product{ID: pid, Name: "Old", Description: "Old Desc", CreatedAt: time.Now()}
	updated := &model.Product{ID: pid, Name: "Updated", Description: "Updated Desc", CreatedAt: time.Now()}

	pr.On("FindByID", ctx, pid.String()).Return(original, nil).Once()
	pr.On("Update", ctx, matchedProduct("Updated", "Updated Desc")).Return(nil).Once()
	pr.On("FindByID", ctx, pid.String()).Return(updated, nil).Once()

	resp, err := uc.UpdateProduct(ctx, pid.String(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Name)
	pr.AssertExpectations(t)
}

func TestProductUsecase_UpdateProduct_NotFound(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pr.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.UpdateProduct(ctx, "x", dto.UpdateProductRequest{Name: "N"})
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	assert.Nil(t, resp)
	pr.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_Success(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pid := uuid.New().String()
	pr.On("FindByID", ctx, pid).Return(&model.Product{ID: uuid.MustParse(pid)}, nil)
	pr.On("Delete", ctx, pid).Return(nil)

	err := uc.DeleteProduct(ctx, pid)
	assert.NoError(t, err)
	pr.AssertExpectations(t)
}

func TestProductUsecase_DeleteProduct_NotFound(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pr.On("FindByID", ctx, "x").Return(nil, nil)

	err := uc.DeleteProduct(ctx, "x")
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	pr.AssertExpectations(t)
}

func TestProductUsecase_CreateVariant_Success(t *testing.T) {
	uc, pr, pvr := setupProductUsecase()
	ctx := context.Background()

	pid := uuid.New()
	req := dto.CreateVariantRequest{Name: "V1", Stock: 10, Price: 99.99}

	pr.On("FindByID", ctx, pid.String()).Return(&model.Product{ID: pid, Name: "P"}, nil)
	pvr.On("Create", ctx, mock.MatchedBy(func(v *model.ProductVariant) bool {
		return v.Name == req.Name && v.Stock == req.Stock && v.Price == req.Price
	})).Return(nil)

	resp, err := uc.CreateVariant(ctx, pid.String(), req)
	assert.NoError(t, err)
	assert.Equal(t, "V1", resp.Name)
	assert.Equal(t, pid.String(), resp.ProductID)
	pr.AssertExpectations(t)
	pvr.AssertExpectations(t)
}

func TestProductUsecase_CreateVariant_ProductNotFound(t *testing.T) {
	uc, pr, _ := setupProductUsecase()
	ctx := context.Background()

	pr.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.CreateVariant(ctx, "x", dto.CreateVariantRequest{Name: "V", Price: 10})
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
	assert.Nil(t, resp)
	pr.AssertExpectations(t)
}

func TestProductUsecase_UpdateVariant_Success(t *testing.T) {
	uc, _, pvr := setupProductUsecase()
	ctx := context.Background()

	vid := uuid.New()
	stock := 20
	req := dto.UpdateVariantRequest{Name: "Updated", Stock: &stock, Price: 199.99}

	original := &model.ProductVariant{ID: vid, Name: "Old", Stock: 10, Price: 99.99}
	pvr.On("FindByID", ctx, vid.String()).Return(original, nil)
	pvr.On("Update", ctx, mock.MatchedBy(func(v *model.ProductVariant) bool {
		return v.Name == "Updated" && v.Stock == 20 && v.Price == 199.99
	})).Return(nil)

	resp, err := uc.UpdateVariant(ctx, vid.String(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated", resp.Name)
	assert.Equal(t, 20, resp.Stock)
	assert.Equal(t, 199.99, resp.Price)
	pvr.AssertExpectations(t)
}

func TestProductUsecase_UpdateVariant_NotFound(t *testing.T) {
	uc, _, pvr := setupProductUsecase()
	ctx := context.Background()

	pvr.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.UpdateVariant(ctx, "x", dto.UpdateVariantRequest{Name: "N"})
	assert.Error(t, err)
	assert.Equal(t, "variant not found", err.Error())
	assert.Nil(t, resp)
	pvr.AssertExpectations(t)
}

func matchedProduct(name, desc string) interface{} {
	return mock.MatchedBy(func(p *model.Product) bool {
		return p.Name == name && p.Description == desc
	})
}
