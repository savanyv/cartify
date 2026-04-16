package model

import "github.com/google/uuid"

type OrderItem struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID          uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	ProductVariantID uuid.UUID `json:"product_variant_id" gorm:"type:uuid;not null"`
	Qty              int       `json:"qty" gorm:"not null"`
	Price            float64   `json:"price" gorm:"not null"`

	Order          Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
}
