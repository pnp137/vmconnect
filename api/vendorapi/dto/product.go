package dto

import "time"

// AddProductRequest represents the request to add a new product
type AddProductRequest struct {
	Name         string  `json:"name" validate:"required,min=1,max=255"`
	Description  string  `json:"description" validate:"max=1000"`
	Price        float64 `json:"price" validate:"required,min=0"`
	CategoryID   uint    `json:"category_id,omitempty" validate:"omitempty"`
	CategoryName string  `json:"category_name" validate:"required,min=1,max=100"`
	SKU          int     `json:"sku"`
	IsActive     bool    `json:"is_active"`
	IsFeatured   bool    `json:"is_featured"`
	ImageURL     string  `json:"image_url" validate:"max=500"`
	Stock        int     `json:"stock" validate:"min=0"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	SKU         *string  `json:"sku,omitempty" validate:"omitempty,max=50"`
	IsActive    *bool    `json:"is_active,omitempty"`
	IsFeatured  *bool    `json:"is_featured,omitempty"`
	ImageURL    *string  `json:"image_url,omitempty" validate:"omitempty,max=500"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,min=0"`
}

// ProductResponse represents the response for product operations
type ProductResponse struct {
	ID           uint      `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float64   `json:"price"`
	CategoryID   uint      `json:"category_id"`
	CategoryName string    `json:"category_name"`
	SKU          string    `json:"sku"`
	IsActive     bool      `json:"is_active"`
	IsFeatured   bool      `json:"is_featured"`
	ImageURL     string    `json:"image_url"`
	Stock        int       `json:"stock"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AddProductResponse represents the response for adding a product
type AddProductResponse struct {
	Product ProductResponse `json:"product"`
	Message string          `json:"message"`
}

// UpdateProductResponse represents the response for updating a product
type UpdateProductResponse struct {
	Product ProductResponse `json:"product"`
	Message string          `json:"message"`
}

// DeleteProductResponse represents the response for deleting a product
type DeleteProductResponse struct {
	ProductID uint   `json:"product_id"`
	Message   string `json:"message"`
}

// GetProductsResponse represents the response for getting products
type GetProductsResponse struct {
	Products []ProductResponse `json:"products"`
	Count    int               `json:"count"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Total    int               `json:"total"`
}

// CategoryInfo represents category information
type CategoryInfo struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// GetCategoriesResponse represents the response for getting vendor categories
type GetCategoriesResponse struct {
	Categories []CategoryInfo `json:"categories"`
	Count      int            `json:"count"`
}

// ProductQueryParam represents query parameters for product listing
type ProductQueryParam struct {
	CategoryID uint   `query:"category_id" validate:"omitempty"`
	IsActive   *bool  `query:"is_active" validate:"omitempty"`
	IsFeatured *bool  `query:"is_featured" validate:"omitempty"`
	Search     string `query:"search" validate:"omitempty,max=100"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
	Offset     int    `query:"offset" validate:"min=0"`
	Ordering   string `query:"ordering" validate:"omitempty,oneof=name price created_at updated_at"`
}
