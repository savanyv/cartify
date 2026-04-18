package model

import (
	"github.com/google/uuid"
)

type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Stock     int       `json:"stock" gorm:"default:0"`
	Price     float64   `json:"price" gorm:"not null"`

	Product Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (pv *ProductVariant) HasStock(qty int) bool {
	return pv.Stock >= qty
}

func (pv *ProductVariant) ReduceStock(qty int) {
	pv.Stock -= qty
}

func (pv *ProductVariant) AddStock(qty int) {
	pv.Stock += qty
}