package tests

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/model"
	"github.com/savanyv/cartify/internal/tests/mocks"
	"github.com/savanyv/cartify/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupOrderUsecaseDB(t *testing.T) (*usecase.OrderUsecase, *mocks.OrderRepositoryMock, *mocks.CartRepositoryMock, *mocks.ProductVariantRepositoryMock, sqlmock.Sqlmock) {
	or := new(mocks.OrderRepositoryMock)
	cr := new(mocks.CartRepositoryMock)
	pvr := new(mocks.ProductVariantRepositoryMock)

	sqlDB, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return usecase.NewOrderUsecase(gormDB, or, cr, pvr), or, cr, pvr, sqlMock
}

func setupOrderUsecase() (*usecase.OrderUsecase, *mocks.OrderRepositoryMock, *mocks.CartRepositoryMock, *mocks.ProductVariantRepositoryMock) {
	or := new(mocks.OrderRepositoryMock)
	cr := new(mocks.CartRepositoryMock)
	pvr := new(mocks.ProductVariantRepositoryMock)
	db := &gorm.DB{}
	return usecase.NewOrderUsecase(db, or, cr, pvr), or, cr, pvr
}

func TestOrderUsecase_CreateOrder_Success(t *testing.T) {
	uc, or, cr, pvr, sqlMock := setupOrderUsecaseDB(t)
	ctx := context.Background()

	variantID := uuid.New()
	productID := uuid.New()
	cartID := uuid.New()
	itemID := uuid.New()

	cart := &model.Cart{
		ID: cartID,
		Items: []model.CartItem{
			{
				ID:               itemID,
				ProductVariantID: variantID,
				Quantity:         2,
				Price:            50.0,
				ProductVariant: model.ProductVariant{
					ID:    variantID,
					Name:  "V1",
					Stock: 10,
					Price: 50.0,
					Product: model.Product{
						ID:   productID,
						Name: "P1",
					},
				},
			},
		},
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	cr.On("GetCartWithItems", ctx, "user1").Return(cart, nil)
	pvr.On("UpdateStock", ctx, variantID.String(), 8).Return(nil)
	or.On("Create", ctx, mock.MatchedBy(func(o *model.Order) bool {
		return o.Status == model.OrderStatusPending && o.TotalPrice == 100.0
	})).Return(nil)
	cr.On("ClearCart", ctx, cartID.String()).Return(nil)

	resp, err := uc.CreateOrder(ctx, "user1")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "pending", resp.Status)
	assert.Equal(t, 100.0, resp.TotalPrice)
	assert.Len(t, resp.Items, 1)
	cr.AssertExpectations(t)
	or.AssertExpectations(t)
	pvr.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestOrderUsecase_CreateOrder_EmptyCart(t *testing.T) {
	uc, _, cr, _, _ := setupOrderUsecaseDB(t)
	ctx := context.Background()

	cr.On("GetCartWithItems", ctx, "u").Return(&model.Cart{ID: uuid.New(), Items: []model.CartItem{}}, nil)

	resp, err := uc.CreateOrder(ctx, "u")
	assert.Error(t, err)
	assert.Equal(t, "cart is empty", err.Error())
	assert.Nil(t, resp)
	cr.AssertExpectations(t)
}

func TestOrderUsecase_CreateOrder_InsufficientStock(t *testing.T) {
	uc, _, cr, _, sqlMock := setupOrderUsecaseDB(t)
	ctx := context.Background()

	variantID := uuid.New()
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	cr.On("GetCartWithItems", ctx, "u").Return(&model.Cart{
		ID: uuid.New(),
		Items: []model.CartItem{
			{
				ID:               uuid.New(),
				ProductVariantID: variantID,
				Quantity:         10,
				Price:            10,
				ProductVariant: model.ProductVariant{
					ID:    variantID,
					Name:  "V1",
					Stock: 1,
				},
			},
		},
	}, nil)

	resp, err := uc.CreateOrder(ctx, "u")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient stock")
	assert.Nil(t, resp)
	cr.AssertExpectations(t)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestOrderUsecase_GetuserOrders_Success(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	now := time.Now()
	orders := []model.Order{
		{
			ID: uuid.New(), UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Status: model.OrderStatusPending, TotalPrice: 100, CreatedAt: now,
			Items: []model.OrderItem{},
		},
	}

	or.On("FindByUserID", ctx, "user1", 1, 10, "", "created_at", "desc").Return(orders, int64(1), nil)

	resp, total, err := uc.GetuserOrders(ctx, "user1", 1, 10, "", "created_at", "desc")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, resp, 1)
	assert.Equal(t, "pending", resp[0].Status)
	assert.Equal(t, 100.0, resp[0].TotalPrice)
	or.AssertExpectations(t)
}

func TestOrderUsecase_GetOrderByID_Success(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()

	order := &model.Order{
		ID: orderID, UserID: userID,
		Status: model.OrderStatusPending, TotalPrice: 100, CreatedAt: time.Now(),
		Items: []model.OrderItem{},
	}

	or.On("FindByID", ctx, orderID.String()).Return(order, nil)

	resp, err := uc.GetOrderByID(ctx, userID.String(), orderID.String())
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, orderID.String(), resp.ID)
	assert.Equal(t, "pending", resp.Status)
	or.AssertExpectations(t)
}

func TestOrderUsecase_GetOrderByID_NotFound(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	or.On("FindByID", ctx, "x").Return(nil, nil)

	resp, err := uc.GetOrderByID(ctx, "u", "x")
	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
	assert.Nil(t, resp)
	or.AssertExpectations(t)
}

func TestOrderUsecase_GetOrderByID_WrongUser(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	order := &model.Order{
		ID: uuid.New(), UserID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		Status: model.OrderStatusPending, TotalPrice: 100, CreatedAt: time.Now(),
	}

	or.On("FindByID", ctx, order.ID.String()).Return(order, nil)

	resp, err := uc.GetOrderByID(ctx, uuid.New().String(), order.ID.String())
	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
	assert.Nil(t, resp)
	or.AssertExpectations(t)
}

func TestOrderUsecase_GetAllOrders_Success(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	now := time.Now()
	orders := []model.Order{
		{
			ID: uuid.New(), UserID: uuid.New(),
			Status: model.OrderStatusPending, TotalPrice: 200, CreatedAt: now,
			Items: []model.OrderItem{},
			User:  model.User{ID: uuid.New(), Name: "U", Email: "u@u.com"},
		},
	}

	or.On("FindAll", ctx, 1, 10, "", "created_at", "desc").Return(orders, int64(1), nil)

	resp, total, err := uc.GetAllOrders(ctx, 1, 10, "", "created_at", "desc")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, resp, 1)
	assert.Equal(t, "pending", resp[0].Status)
	assert.Equal(t, 200.0, resp[0].TotalPrice)
	or.AssertExpectations(t)
}

func TestOrderUsecase_UpdateOrderStatus_Success(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	orderID := uuid.New().String()
	or.On("FindByID", ctx, orderID).Return(&model.Order{ID: uuid.MustParse(orderID), Status: model.OrderStatusPending}, nil)
	or.On("UpdateStatus", ctx, orderID, model.OrderStatusPaid).Return(nil)

	err := uc.UpdateOrderStatus(ctx, orderID, model.OrderStatusPaid)
	assert.NoError(t, err)
	or.AssertExpectations(t)
}

func TestOrderUsecase_UpdateOrderStatus_NotFound(t *testing.T) {
	uc, or, _, _ := setupOrderUsecase()
	ctx := context.Background()

	or.On("FindByID", ctx, "x").Return(nil, nil)

	err := uc.UpdateOrderStatus(ctx, "x", model.OrderStatusShipped)
	assert.Error(t, err)
	assert.Equal(t, "order not found", err.Error())
	or.AssertExpectations(t)
}
