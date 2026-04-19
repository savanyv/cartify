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

func (r *productRepository) FindAllWithPagination(ctx context.Context, page, limit int, search, sort, order string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Product{}).Preload("Variants")

	// Apply search
	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(sort + " " + order)

	// Apply pagination
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, ID string) error {
	return r.db.WithContext(ctx).Where("id = ?", ID).Delete(&model.Product{}).Error
}