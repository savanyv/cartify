package dto

type CreateProductRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"max=100"`
}

type UpdateProductRequest struct {
	Name string `json:"name" validate:"omitempty,min=3,max=100"`
	Description string `json:"description" validate:"omitempty,max=100"`
}

type CreateVariantRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	Stock int `json:"stock" validate:"min=0"`
	Price float64 `json:"price" validate:"required,gt=0"`
}

type UpdateVariantRequest struct {
	Name string `json:"name" validate:"omitempty,min=1,max=100"`
	Stock *int `json:"stock" validate:"omitempty,min=0"`
	Price float64 `json:"price" validate:"omitempty,gt=0"`
}

type VariantResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Stock int `json:"stock"`
	Price float64 `json:"price"`
	ProductID string `json:"product_id"`
}

type ProductResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt string `json:"created_at"`
	Variants []VariantResponse `json:"variants"`
}