package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/savanyv/cartify/internal/model"
	"gorm.io/gorm"
)

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) model.CartRepository {
	return &cartRepository{
		db: db,
	}
}

func (r *cartRepository) GetOrCreateCart(ctx context.Context, userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new cart
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	cart = model.Cart{UserID: parsedUserID}
	if err := r.db.WithContext(ctx).Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepository) AddItem(ctx context.Context, cartID string, variantID string, price float64, qty int) error {
	parsedCartID, err := uuid.Parse(cartID)
	if err != nil {
		return err
	}
	parsedVariantID, err := uuid.Parse(variantID)
	if err != nil {
		return err
	}

	// Check if item already exists
	var existingItem model.CartItem
	err = r.db.WithContext(ctx).Where("cart_id = ? AND product_variant_id = ?", parsedCartID, parsedVariantID).First(&existingItem).Error
	if err == nil {
		// Update quantity
		newQty := existingItem.Quantity + qty
		return r.db.WithContext(ctx).Model(&existingItem).Update("quantity", newQty).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create new item
	item := model.CartItem{
		CartID:           parsedCartID,
		ProductVariantID: parsedVariantID,
		Quantity:         qty,
		Price:            price,
	}
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *cartRepository) GetCartWithItems(ctx context.Context, userID string) (*model.Cart, error) {
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	var cart model.Cart
	err = r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.ProductVariant").
		Preload("Items.ProductVariant.Product").
		Where("user_id = ?", parsedUserID).
		First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.Cart{UserID: parsedUserID, Items: []model.CartItem{}}, nil
	}
	return &cart, err
}

func (r *cartRepository) GetCartItem(ctx context.Context, cartID string, variantID string) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_variant_id = ?", cartID, variantID).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, cartItemID string, qty int) error {
	parsedID, err := uuid.Parse(cartItemID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(&model.CartItem{}).Where("id = ?", parsedID).Update("quantity", qty).Error
}

func (r *cartRepository) RemoveItem(ctx context.Context, cartItemID string) error {
	parsedID, err := uuid.Parse(cartItemID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(&model.CartItem{}, parsedID).Error
}

func (r *cartRepository) ClearCart(ctx context.Context, cartID string) error {
	parsedID, err := uuid.Parse(cartID)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("cart_id = ?", parsedID).Delete(&model.CartItem{}).Error
}