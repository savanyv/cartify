package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusPaid, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

type Order struct {
	ID         uuid.UUID   `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID     uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	Status     OrderStatus `json:"status" gorm:"type:varchar(50);not null;default:pending"`
	TotalPrice float64     `json:"total_price" gorm:"not null"`
	CreatedAt  time.Time   `json:"created_at" gorm:"autoCreateTime"`

	User  User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (o *Order) CanBeCancelled() bool {
	return o.Status == OrderStatusPending
}

func (o *Order) CanBePaid() bool {
	return o.Status == OrderStatusPending
}

