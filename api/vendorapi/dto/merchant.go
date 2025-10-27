package dto

import "time"

// MerchantInfo represents merchant information for vendor
type MerchantInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ShopName  string    `json:"shop_name"`
	Address   string    `json:"address"`
	Pincode   string    `json:"pincode"`
	CreatedAt time.Time `json:"created_at"`
}

// GetMerchantsResponse represents the response for getting vendor's merchants
type GetMerchantsResponse struct {
	Merchants []MerchantInfo `json:"merchants"`
	Count     int            `json:"count"`
}

// ProductVisibilityInfo represents product visibility for a merchant
type ProductVisibilityInfo struct {
	ProductID    uint    `json:"product_id"`
	ProductName  string  `json:"product_name"`
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Price        float64 `json:"price"`
	Visible      bool    `json:"visible"`
}

// GetProductVisibilityResponse represents the response for getting product visibility
type GetProductVisibilityResponse struct {
	Products []ProductVisibilityInfo `json:"products"`
	Count    int                     `json:"count"`
}

// UpdateProductVisibilityRequest represents the request to update product visibility
type UpdateProductVisibilityRequest struct {
	ProductID uint `json:"product_id" validate:"required"`
	Visible   bool `json:"visible"`
}

// UpdateProductVisibilityResponse represents the response for updating product visibility
type UpdateProductVisibilityResponse struct {
	UpdatedCount int    `json:"updated_count"`
	Message      string `json:"message"`
}
