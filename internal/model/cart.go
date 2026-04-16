package model

import (
	"github.com/google/uuid"
)

type Cart struct {
	ID     uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null;unique"`

	User  User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items []CartItem `json:"items,omitempty" gorm:"foreignKey:CartID"`
}

func (c *Cart) GetTotalPrice() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.ProductVariant.Price * float64(item.Quantity)
	}
	return total
}