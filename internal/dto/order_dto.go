package dto

type UpdateOrderStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending paid shipped delivered cancelled"`
}

type OrderItemResponse struct {
	ID               string  `json:"id"`
	ProductVariantID string  `json:"product_variant_id"`
	ProductName      string  `json:"product_name"`
	VariantName      string  `json:"variant_name"`
	Qty              int     `json:"qty"`
	Price            float64 `json:"price"`
	SubTotal         float64 `json:"sub_total"`
}

type OrderResponse struct {
	ID         string              `json:"id"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"total_price"`
	CreatedAt  string              `json:"created_at"`
	Items      []OrderItemResponse `json:"items"`
}

type AdminOrderResponse struct {
	ID         string              `json:"id"`
	UserID     string              `json:"user_id"`
	UserName   string              `json:"user_name"`
	UserEmail  string              `json:"user_email"`
	Status     string              `json:"status"`
	TotalPrice float64             `json:"total_price"`
	CreatedAt  string              `json:"created_at"`
	Items      []OrderItemResponse `json:"items"`
}
