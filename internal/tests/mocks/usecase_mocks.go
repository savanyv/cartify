package mocks

import (
	"context"

	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.RegisterResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AuthUsecaseMock) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.LoginResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AuthUsecaseMock) GetUserByID(ctx context.Context, id string) (*dto.UserResponse, error) {
	args := m.Called(ctx, id)
	if resp, ok := args.Get(0).(*dto.UserResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AuthUsecaseMock) ChangePassword(ctx context.Context, userID string, req dto.ChangePasswordRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *AuthUsecaseMock) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	args := m.Called(ctx, refreshToken)
	if resp, ok := args.Get(0).(*dto.RefreshTokenResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *AuthUsecaseMock) Logout(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type ProductUsecaseMock struct {
	mock.Mock
}

func (m *ProductUsecaseMock) GetAllProductsWithPagination(ctx context.Context, page, limit int, search, sort, order string) ([]dto.ProductResponse, int64, error) {
	args := m.Called(ctx, page, limit, search, sort, order)
	if resp, ok := args.Get(0).([]dto.ProductResponse); ok {
		return resp, args.Get(1).(int64), args.Error(2)
	}
	return nil, int64(0), args.Error(2)
}

func (m *ProductUsecaseMock) GetProductByID(ctx context.Context, id string) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id)
	if resp, ok := args.Get(0).(*dto.ProductResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUsecaseMock) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	args := m.Called(ctx, req)
	if resp, ok := args.Get(0).(*dto.ProductResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUsecaseMock) UpdateProduct(ctx context.Context, id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	args := m.Called(ctx, id, req)
	if resp, ok := args.Get(0).(*dto.ProductResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUsecaseMock) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductUsecaseMock) CreateVariant(ctx context.Context, productID string, req dto.CreateVariantRequest) (*dto.VariantResponse, error) {
	args := m.Called(ctx, productID, req)
	if resp, ok := args.Get(0).(*dto.VariantResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductUsecaseMock) UpdateVariant(ctx context.Context, id string, req dto.UpdateVariantRequest) (*dto.VariantResponse, error) {
	args := m.Called(ctx, id, req)
	if resp, ok := args.Get(0).(*dto.VariantResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

type CartUsecaseMock struct {
	mock.Mock
}

func (m *CartUsecaseMock) GetCart(ctx context.Context, userID string) (*dto.CartResponse, error) {
	args := m.Called(ctx, userID)
	if resp, ok := args.Get(0).(*dto.CartResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartUsecaseMock) AddToCart(ctx context.Context, userID string, req dto.AddToCartRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *CartUsecaseMock) UpdateCartItem(ctx context.Context, userID string, itemID string, req dto.UpdateCartItemRequest) error {
	args := m.Called(ctx, userID, itemID, req)
	return args.Error(0)
}

func (m *CartUsecaseMock) RemoveFromCart(ctx context.Context, userID string, itemID string) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}

func (m *CartUsecaseMock) ClearCart(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type OrderUsecaseMock struct {
	mock.Mock
}

func (m *OrderUsecaseMock) CreateOrder(ctx context.Context, userID string) (*dto.OrderResponse, error) {
	args := m.Called(ctx, userID)
	if resp, ok := args.Get(0).(*dto.OrderResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrderUsecaseMock) GetuserOrders(ctx context.Context, userID string, page, limit int, search, sort, order string) ([]dto.OrderResponse, int64, error) {
	args := m.Called(ctx, userID, page, limit, search, sort, order)
	if resp, ok := args.Get(0).([]dto.OrderResponse); ok {
		return resp, args.Get(1).(int64), args.Error(2)
	}
	return nil, int64(0), args.Error(2)
}

func (m *OrderUsecaseMock) GetOrderByID(ctx context.Context, userID string, orderID string) (*dto.OrderResponse, error) {
	args := m.Called(ctx, userID, orderID)
	if resp, ok := args.Get(0).(*dto.OrderResponse); ok {
		return resp, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrderUsecaseMock) GetAllOrders(ctx context.Context, page, limit int, search, sort, order string) ([]dto.AdminOrderResponse, int64, error) {
	args := m.Called(ctx, page, limit, search, sort, order)
	if resp, ok := args.Get(0).([]dto.AdminOrderResponse); ok {
		return resp, args.Get(1).(int64), args.Error(2)
	}
	return nil, int64(0), args.Error(2)
}

func (m *OrderUsecaseMock) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}
