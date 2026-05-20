package mocks

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/utils/helpers"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if user != nil && user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if u, ok := args.Get(0).(*model.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	if u, ok := args.Get(0).(*model.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, ID string) (*model.User, error) {
	args := m.Called(ctx, ID)
	if u, ok := args.Get(0).(*model.User); ok {
		return u, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *UserRepositoryMock) UpdateTokenVersion(ctx context.Context, ID string, version int) error {
	args := m.Called(ctx, ID, version)
	return args.Error(0)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Create(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) FindByID(ctx context.Context, ID string) (*model.Product, error) {
	args := m.Called(ctx, ID)
	if p, ok := args.Get(0).(*model.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductRepositoryMock) FindAll(ctx context.Context) ([]model.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindAllWithPagination(ctx context.Context, page, limit int, search, sort, order string) ([]model.Product, int64, error) {
	args := m.Called(ctx, page, limit, search, sort, order)
	return args.Get(0).([]model.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *model.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, ID string) error {
	args := m.Called(ctx, ID)
	return args.Error(0)
}

type ProductVariantRepositoryMock struct {
	mock.Mock
}

func (m *ProductVariantRepositoryMock) Create(ctx context.Context, variant *model.ProductVariant) error {
	args := m.Called(ctx, variant)
	if variant != nil && variant.ID == uuid.Nil {
		variant.ID = uuid.New()
	}
	return args.Error(0)
}

func (m *ProductVariantRepositoryMock) FindByID(ctx context.Context, ID string) (*model.ProductVariant, error) {
	args := m.Called(ctx, ID)
	if v, ok := args.Get(0).(*model.ProductVariant); ok {
		return v, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductVariantRepositoryMock) FindByProductID(ctx context.Context, productID string) ([]model.ProductVariant, error) {
	args := m.Called(ctx, productID)
	return args.Get(0).([]model.ProductVariant), args.Error(1)
}

func (m *ProductVariantRepositoryMock) Update(ctx context.Context, variant *model.ProductVariant) error {
	args := m.Called(ctx, variant)
	return args.Error(0)
}

func (m *ProductVariantRepositoryMock) UpdateStock(ctx context.Context, ID string, stock int) error {
	args := m.Called(ctx, ID, stock)
	return args.Error(0)
}

type CartRepositoryMock struct {
	mock.Mock
}

func (m *CartRepositoryMock) GetOrCreateCart(ctx context.Context, userID string) (*model.Cart, error) {
	args := m.Called(ctx, userID)
	if c, ok := args.Get(0).(*model.Cart); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartRepositoryMock) AddItem(ctx context.Context, cartID string, variantID string, price float64, qty int) error {
	args := m.Called(ctx, cartID, variantID, price, qty)
	return args.Error(0)
}

func (m *CartRepositoryMock) GetCartWithItems(ctx context.Context, userID string) (*model.Cart, error) {
	args := m.Called(ctx, userID)
	if c, ok := args.Get(0).(*model.Cart); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartRepositoryMock) GetCartItem(ctx context.Context, cartID string, variantID string) (*model.CartItem, error) {
	args := m.Called(ctx, cartID, variantID)
	if ci, ok := args.Get(0).(*model.CartItem); ok {
		return ci, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *CartRepositoryMock) UpdateItemQuantity(ctx context.Context, cartItemID string, qty int) error {
	args := m.Called(ctx, cartItemID, qty)
	return args.Error(0)
}

func (m *CartRepositoryMock) RemoveItem(ctx context.Context, cartItemID string) error {
	args := m.Called(ctx, cartItemID)
	return args.Error(0)
}

func (m *CartRepositoryMock) ClearCart(ctx context.Context, cartID string) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) Create(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) FindByID(ctx context.Context, ID string) (*model.Order, error) {
	args := m.Called(ctx, ID)
	if o, ok := args.Get(0).(*model.Order); ok {
		return o, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *OrderRepositoryMock) FindByUserID(ctx context.Context, userID string, page, limit int, search, sort, order string) ([]model.Order, int64, error) {
	args := m.Called(ctx, userID, page, limit, search, sort, order)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) FindAll(ctx context.Context, page, limit int, search, sort, order string) ([]model.Order, int64, error) {
	args := m.Called(ctx, page, limit, search, sort, order)
	return args.Get(0).([]model.Order), args.Get(1).(int64), args.Error(2)
}

func (m *OrderRepositoryMock) UpdateStatus(ctx context.Context, ID string, status model.OrderStatus) error {
	args := m.Called(ctx, ID, status)
	return args.Error(0)
}

type BcryptServiceMock struct {
	mock.Mock
}

func (m *BcryptServiceMock) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *BcryptServiceMock) ComparePassword(hashedPassword, password string) bool {
	args := m.Called(hashedPassword, password)
	return args.Bool(0)
}

type JWTServiceMock struct {
	mock.Mock
}

func (m *JWTServiceMock) GenerateAccessToken(userID, username, email, role string, tokenVersion int) (string, error) {
	args := m.Called(userID, username, email, role, tokenVersion)
	return args.String(0), args.Error(1)
}

func (m *JWTServiceMock) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *JWTServiceMock) ValidateToken(tokenString string) (*helpers.JWTClaims, error) {
	args := m.Called(tokenString)
	if c, ok := args.Get(0).(*helpers.JWTClaims); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *JWTServiceMock) ValidateRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	args := m.Called(tokenString)
	if c, ok := args.Get(0).(*jwt.RegisteredClaims); ok {
		return c, args.Error(1)
	}
	return nil, args.Error(1)
}
