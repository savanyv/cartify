package dto

type AddToCartRequest struct {
	ProductVariantID string `json:"product_variant_id"`
	Quantity         int    `json:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

type CartItemResponse struct {
	ID string `json:"id"`
	ProductVariantID string `json:"product_variant_id"`
	ProductName string `json:"product_name"`
	VariantName string `json:"variant_name"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
	SubTotal float64 `json:"sub_total"`
}

type CartResponse struct {
	ID string `json:"id"`
	Items []CartItemResponse `json:"items"`
	TotalPrice float64 `json:"total_price"`
	ItemCount int `json:"item_count"`
}