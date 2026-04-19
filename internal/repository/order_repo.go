package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/model"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) model.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *model.Order) error {
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderRepository) FindByID(ctx context.Context, ID string) (*model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Product").
		Where("id = ?", ID).
		First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepository) FindByUserID(ctx context.Context, userID string, page, limit int, search, sort, order string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&model.Order{}).Where("user_id = ?", parsedUserID)

	if search != "" {
		query = query.Where("status ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(sort + " " + order)

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Product").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) FindAll(ctx context.Context, page, limit int, search, sort, order string) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Order{})

	if search != "" {
		query = query.Joins("JOIN users ON users.id = orders.user_id").
			Where("orders.status ILIKE ? OR users.email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(sort + " " + order)

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Preload("User").
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Product").
		Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) UpdateStatus(ctx context.Context, ID string, status model.OrderStatus) error {
	parsedID, err := uuid.Parse(ID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", parsedID).Update("status", status).Error
}