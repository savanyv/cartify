package model

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByID(ctx context.Context, ID string) (*User, error)
	UpdateTokenVersion(ctx context.Context, ID string, version int) error
	Update(ctx context.Context, user *User) error
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, ID string) (*Product, error)
	FindAll(ctx context.Context) ([]Product, error)
	FindAllWithPagination(ctx context.Context, page, limit int, search, sort, order string) ([]Product, int64, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, ID string) error
}

type ProductVariantRepository interface {
	Create(ctx context.Context, variant *ProductVariant) error
	FindByID(ctx context.Context, ID string) (*ProductVariant, error)
	FindByProductID(ctx context.Context, productID string) ([]ProductVariant, error)
	Update(ctx context.Context, variant *ProductVariant) error
	UpdateStock(ctx context.Context, ID string, stock int) error
}

type CartRepository interface {
	GetOrCreateCart(ctx context.Context, userID string) (*Cart, error)
	AddItem(ctx context.Context, cartID string, variantID string, price float64, qty int) error
	GetCartWithItems(ctx context.Context, userID string) (*Cart, error)
	GetCartItem(ctx context.Context, cartID string, variantID string) (*CartItem, error)
	UpdateItemQuantity(ctx context.Context, cartItemID string, qty int) error
	RemoveItem(ctx context.Context, cartItemID string) error
	ClearCart(ctx context.Context, cartID string) error
}