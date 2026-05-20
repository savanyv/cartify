package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCartUsecase() (*usecase.CartUsecase, *mocks.CartRepositoryMock, *mocks.ProductVariantRepositoryMock) {
	cr := new(mocks.CartRepositoryMock)
	pvr := new(mocks.ProductVariantRepositoryMock)
	return usecase.NewCartUsecase(cr, pvr), cr, pvr
}

func TestCartUsecase_GetCart_Success(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cartID := uuid.New()
	variantID := uuid.New()
	productID := uuid.New()

	cart := &model.Cart{
		ID: cartID,
		Items: []model.CartItem{
			{
				ID:               uuid.New(),
				ProductVariantID: variantID,
				Quantity:         2,
				Price:            50.0,
				ProductVariant: model.ProductVariant{
					ID:   variantID,
					Name: "V1",
					Product: model.Product{
						ID:   productID,
						Name: "P1",
					},
				},
			},
		},
	}

	cr.On("GetCartWithItems", ctx, "user1").Return(cart, nil)

	resp, err := uc.GetCart(ctx, "user1")
	assert.NoError(t, err)
	assert.Equal(t, cartID.String(), resp.ID)
	assert.Len(t, resp.Items, 1)
	assert.Equal(t, "P1", resp.Items[0].ProductName)
	assert.Equal(t, "V1", resp.Items[0].VariantName)
	assert.Equal(t, 100.0, resp.TotalPrice)
	assert.Equal(t, 1, resp.ItemCount)
	cr.AssertExpectations(t)
}

func TestCartUsecase_GetCart_Error(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cr.On("GetCartWithItems", ctx, "user1").Return(nil, errors.New("db error"))

	resp, err := uc.GetCart(ctx, "user1")
	assert.Error(t, err)
	assert.Nil(t, resp)
	cr.AssertExpectations(t)
}

func TestCartUsecase_AddToCart_Success(t *testing.T) {
	uc, cr, pvr := setupCartUsecase()
	ctx := context.Background()

	variantID := uuid.New().String()
	req := dto.AddToCartRequest{ProductVariantID: variantID, Quantity: 2}

	pvr.On("FindByID", ctx, variantID).Return(&model.ProductVariant{ID: uuid.MustParse(variantID), Stock: 10, Price: 25.0}, nil)
	cr.On("GetOrCreateCart", ctx, "user1").Return(&model.Cart{ID: uuid.New()}, nil)
	cr.On("AddItem", ctx, mock.Anything, variantID, 25.0, 2).Return(nil)

	err := uc.AddToCart(ctx, "user1", req)
	assert.NoError(t, err)
	pvr.AssertExpectations(t)
	cr.AssertExpectations(t)
}

func TestCartUsecase_AddToCart_InsufficientStock(t *testing.T) {
	uc, _, pvr := setupCartUsecase()
	ctx := context.Background()

	variantID := uuid.New().String()
	pvr.On("FindByID", ctx, variantID).Return(&model.ProductVariant{ID: uuid.MustParse(variantID), Stock: 1}, nil)

	err := uc.AddToCart(ctx, "u", dto.AddToCartRequest{ProductVariantID: variantID, Quantity: 5})
	assert.Error(t, err)
	assert.Equal(t, "insufficient stock", err.Error())
	pvr.AssertExpectations(t)
}

func TestCartUsecase_AddToCart_VariantNotFound(t *testing.T) {
	uc, _, pvr := setupCartUsecase()
	ctx := context.Background()

	pvr.On("FindByID", ctx, "x").Return(nil, nil)

	err := uc.AddToCart(ctx, "u", dto.AddToCartRequest{ProductVariantID: "x", Quantity: 1})
	assert.Error(t, err)
	assert.Equal(t, "product variant not found", err.Error())
	pvr.AssertExpectations(t)
}

func TestCartUsecase_UpdateCartItem_Success(t *testing.T) {
	uc, cr, pvr := setupCartUsecase()
	ctx := context.Background()

	itemID := uuid.New().String()
	variantID := uuid.New()

	cr.On("GetCartWithItems", ctx, "user1").Return(&model.Cart{
		Items: []model.CartItem{
			{
				ID:               uuid.MustParse(itemID),
				ProductVariantID: variantID,
				Quantity:         1,
				Price:            10.0,
			},
		},
	}, nil)
	pvr.On("FindByID", ctx, variantID.String()).Return(&model.ProductVariant{ID: variantID, Stock: 10}, nil)
	cr.On("UpdateItemQuantity", ctx, itemID, 5).Return(nil)

	err := uc.UpdateCartItem(ctx, "user1", itemID, dto.UpdateCartItemRequest{Quantity: 5})
	assert.NoError(t, err)
	cr.AssertExpectations(t)
	pvr.AssertExpectations(t)
}

func TestCartUsecase_UpdateCartItem_ItemNotFound(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cr.On("GetCartWithItems", ctx, "u").Return(&model.Cart{Items: []model.CartItem{}}, nil)

	err := uc.UpdateCartItem(ctx, "u", "x", dto.UpdateCartItemRequest{Quantity: 1})
	assert.Error(t, err)
	assert.Equal(t, "item not found in cart", err.Error())
	cr.AssertExpectations(t)
}

func TestCartUsecase_UpdateCartItem_InsufficientStock(t *testing.T) {
	uc, cr, pvr := setupCartUsecase()
	ctx := context.Background()

	itemID := uuid.New()
	variantID := uuid.New()

	cr.On("GetCartWithItems", ctx, "u").Return(&model.Cart{
		Items: []model.CartItem{
			{ID: itemID, ProductVariantID: variantID, Quantity: 1, Price: 10},
		},
	}, nil)
	pvr.On("FindByID", ctx, variantID.String()).Return(&model.ProductVariant{ID: variantID, Stock: 1}, nil)

	err := uc.UpdateCartItem(ctx, "u", itemID.String(), dto.UpdateCartItemRequest{Quantity: 10})
	assert.Error(t, err)
	assert.Equal(t, "insufficient stock", err.Error())
	cr.AssertExpectations(t)
	pvr.AssertExpectations(t)
}

func TestCartUsecase_RemoveFromCart_Success(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	itemID := uuid.New()

	cr.On("GetCartWithItems", ctx, "user1").Return(&model.Cart{
		Items: []model.CartItem{
			{ID: itemID, ProductVariantID: uuid.New(), Quantity: 1, Price: 10},
		},
	}, nil)
	cr.On("RemoveItem", ctx, itemID.String()).Return(nil)

	err := uc.RemoveFromCart(ctx, "user1", itemID.String())
	assert.NoError(t, err)
	cr.AssertExpectations(t)
}

func TestCartUsecase_RemoveFromCart_ItemNotFound(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cr.On("GetCartWithItems", ctx, "u").Return(&model.Cart{Items: []model.CartItem{}}, nil)

	err := uc.RemoveFromCart(ctx, "u", "x")
	assert.Error(t, err)
	assert.Equal(t, "item not found in cart", err.Error())
	cr.AssertExpectations(t)
}

func TestCartUsecase_ClearCart_Success(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cartID := uuid.New()
	cr.On("GetCartWithItems", ctx, "user1").Return(&model.Cart{ID: cartID}, nil)
	cr.On("ClearCart", ctx, cartID.String()).Return(nil)

	err := uc.ClearCart(ctx, "user1")
	assert.NoError(t, err)
	cr.AssertExpectations(t)
}

func TestCartUsecase_ClearCart_Error(t *testing.T) {
	uc, cr, _ := setupCartUsecase()
	ctx := context.Background()

	cr.On("GetCartWithItems", ctx, "u").Return(nil, errors.New("db error"))

	err := uc.ClearCart(ctx, "u")
	assert.Error(t, err)
	cr.AssertExpectations(t)
}
