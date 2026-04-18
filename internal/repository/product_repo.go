package repository

import (
	"context"
	"errors"

	"github.com/savanyv/cartify/internal/model"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) model.ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) FindByID(ctx context.Context, ID string) (*model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).Preload("Variants").Where("id = ?", ID).First(&product).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return &product, err
}

func (r *productRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	err := r.db.WithContext(ctx).Preload("Variants").Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Where("id = ?", ID).Delete(&model.Product{}).Error
}