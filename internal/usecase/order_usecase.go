package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/dto"
	"github.com/savanyv/cartify/internal/model"
	"gorm.io/gorm"
)

type OrderUsecase struct {
	db                 *gorm.DB
	orderRepo          model.OrderRepository
	cartRepo           model.CartRepository
	productVariantRepo model.ProductVariantRepository
}

func NewOrderUsecase(db *gorm.DB, or model.OrderRepository, cr model.CartRepository, pvr model.ProductVariantRepository) *OrderUsecase {
	return &OrderUsecase{
		db:                 db,
		orderRepo:          or,
		cartRepo:           cr,
		productVariantRepo: pvr,
	}
}

func (u *OrderUsecase) CreateOrder(ctx context.Context, userID string) (*dto.OrderResponse, error) {
	cart, err := u.cartRepo.GetCartWithItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	tx := u.db.Begin()

	var totalPrice float64
	var orderItems []model.OrderItem

	for _, item := range cart.Items {
		if item.ProductVariant.Stock < item.Quantity {
			tx.Rollback()
			return nil, errors.New("insufficient stock for " + item.ProductVariant.Name)
		}

		subTotal := item.Price * float64(item.Quantity)
		totalPrice += subTotal

		orderItems = append(orderItems, model.OrderItem{
			ProductVariantID: item.ProductVariantID,
			Qty:              item.Quantity,
			Price:            item.Price,
		})

		newStock := item.ProductVariant.Stock - item.Quantity
		if err := u.productVariantRepo.UpdateStock(ctx, item.ProductVariantID.String(), newStock); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	parsedUserID, _ := uuid.Parse(userID)
	order := model.Order{
		UserID:     parsedUserID,
		Status:     model.OrderStatusPending,
		TotalPrice: totalPrice,
		Items:      orderItems,
	}

	if err := u.orderRepo.Create(ctx, &order); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := u.cartRepo.ClearCart(ctx, cart.ID.String()); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return u.toOrderResponse(&order), nil
}

func (u *OrderUsecase) GetuserOrders(ctx context.Context, userID string, page, limit int, search, sort, order string) ([]dto.OrderResponse, int64, error) {
	orders, total, err := u.orderRepo.FindByUserID(ctx, userID, page, limit, search, sort, order)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.OrderResponse
	for _, o := range orders {
		responses = append(responses, *u.toOrderResponse(&o))
	}

	return responses, total, nil
}

func (u *OrderUsecase) GetOrderByID(ctx context.Context, userID string, orderID string) (*dto.OrderResponse, error) {
	order, err := u.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}

	if order.UserID.String() != userID {
		return nil, errors.New("order not found")
	}

	return u.toOrderResponse(order), nil
}

func (u *OrderUsecase) GetAllOrders(ctx context.Context, page, limit int, search, sort, order string) ([]dto.AdminOrderResponse, int64, error) {
	orders, total, err := u.orderRepo.FindAll(ctx, page, limit, search, sort, order)
	if err != nil {
		return nil, 0, err
	}

	var responses []dto.AdminOrderResponse
	for _, o := range orders {
		responses = append(responses, *u.toAdminOrderResponse(&o))
	}
	return responses, total, nil
}

func (u *OrderUsecase) UpdateOrderStatus(ctx context.Context, orderID string, status model.OrderStatus) error {
	order, err := u.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	return u.orderRepo.UpdateStatus(ctx, orderID, status)
}

func (u *OrderUsecase) toOrderResponse(order *model.Order) *dto.OrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range order.Items {
		productName := ""
		variantName := ""
		if item.ProductVariant.Product.ID != uuid.Nil {
			productName = item.ProductVariant.Product.Name
		}
		variantName = item.ProductVariant.Name

		items = append(items, dto.OrderItemResponse{
			ID:               item.ID.String(),
			ProductVariantID: item.ProductVariantID.String(),
			ProductName:      productName,
			VariantName:      variantName,
			Qty:              item.Qty,
			Price:            item.Price,
			SubTotal:         item.Price * float64(item.Qty),
		})
	}

	return &dto.OrderResponse{
		ID:         order.ID.String(),
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		CreatedAt:  order.CreatedAt.Format(time.RFC3339),
		Items: items,
	}
}

func (u *OrderUsecase) toAdminOrderResponse(order *model.Order) *dto.AdminOrderResponse {
	var items []dto.OrderItemResponse
	for _, item := range order.Items {
		productName := ""
		variantName := ""
		if item.ProductVariant.Product.ID != uuid.Nil {
			productName = item.ProductVariant.Product.Name
		}
		variantName = item.ProductVariant.Name

		items = append(items, dto.OrderItemResponse{
			ID: item.ID.String(),
			ProductVariantID: item.ProductVariantID.String(),
			ProductName: productName,
			VariantName: variantName,
			Qty: item.Qty,
			Price: item.Price,
			SubTotal: item.Price * float64(item.Qty),
		})
	}

	userName := ""
	userEmail := ""
	if order.User.ID != uuid.Nil {
		userName = order.User.Name
		userEmail = order.User.Email
	}

	return &dto.AdminOrderResponse{
		ID: order.ID.String(),
		UserID: order.UserID.String(),
		UserName: userName,
		UserEmail: userEmail,
		Status: string(order.Status),
		TotalPrice: order.TotalPrice,
		CreatedAt: order.CreatedAt.Format(time.RFC3339),
		Items: items,
	}
}
