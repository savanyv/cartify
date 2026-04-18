package repository

import (
	"context"
	"errors"

	"github.com/savanyv/cartify/internal/model"
	"gorm.io/gorm"
)

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) model.ProductVariantRepository {
	return &productVariantRepository{
		db: db,
	}
}

func (r *productVariantRepository) Create(ctx context.Context, variant *model.ProductVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *productVariantRepository) FindByID(ctx context.Context, ID string) (*model.ProductVariant, error) {
	var variant model.ProductVariant
	err := r.db.WithContext(ctx).Preload("Product").Where("id = ?", ID).First(&variant).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &variant, nil
}

func (r *productVariantRepository) FindByProductID(ctx context.Context, productID string) ([]model.ProductVariant, error) {
	var variants []model.ProductVariant
	err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&variants).Error

	return variants, err
}

func (r *productVariantRepository) Update(ctx context.Context, variant *model.ProductVariant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *productVariantRepository) UpdateStock(ctx context.Context, ID string, stock int) error {
	return r.db.WithContext(ctx).Model(&model.ProductVariant{}).Where("id = ?", ID).Update("stock", stock).Error
}