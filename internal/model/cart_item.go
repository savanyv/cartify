package model

import (
	"github.com/google/uuid"
)

type CartItem struct {
	ID               uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CartID           uuid.UUID `json:"cart_id" gorm:"type:uuid;not null"`
	ProductVariantID uuid.UUID `json:"product_variant_id" gorm:"type:uuid;not null"`
	Quantity         int       `json:"quantity" gorm:"default:1"`

	Cart           Cart           `json:"cart,omitempty" gorm:"foreignKey:CartID"`
	ProductVariant ProductVariant `json:"product_variant,omitempty" gorm:"foreignKey:ProductVariantID"`
}