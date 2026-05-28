package dto

import "time"

type ProductVariantRequest struct {
	ID       *uint   `json:"id,omitempty"`
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Unit     string  `json:"unit,omitempty" validate:"omitempty,max=50"`
	Quantity float64 `json:"quantity" validate:"omitempty"`
	Price    float64 `json:"price" validate:"required,min=0"`
	MRP      float64 `json:"mrp" validate:"omitempty,min=0"`
	MOQ      float64 `json:"moq" validate:"omitempty,min=0"`
	Stock    int     `json:"stock" validate:"omitempty,min=0"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type ProductVariantResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Unit      string    `json:"unit"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	MRP       float64   `json:"mrp"`
	MOQ       float64   `json:"moq"`
	Stock     int       `json:"stock"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AddProductRequest represents the request to add a new product
type AddProductRequest struct {
	Name         string                  `json:"name" validate:"required,min=1,max=255"`
	Description  string                  `json:"description" validate:"max=1000"`
	CategoryName string                  `json:"category_name" validate:"required,min=1,max=100"`
	IsActive     bool                    `json:"is_active"`
	IsFeatured   bool                    `json:"is_featured"`
	ImageURLs    []string                `json:"image_urls,omitempty"`
	ThumbnailURL string                  `json:"thumbnail_url,omitempty" validate:"omitempty,max=500"`
	Variants     []ProductVariantRequest `json:"variants" validate:"required,min=1,dive"`
}

// UpdateProductRequest represents the request to update a product
type UpdateProductRequest struct {
	Name         *string                  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description  *string                  `json:"description,omitempty" validate:"omitempty,max=1000"`
	CategoryName *string                  `json:"category_name,omitempty" validate:"omitempty,min=1,max=100"`
	IsActive     *bool                    `json:"is_active,omitempty"`
	IsFeatured   *bool                    `json:"is_featured,omitempty"`
	ImageURLs    *[]string                `json:"image_urls,omitempty"`
	ThumbnailURL *string                  `json:"thumbnail_url,omitempty" validate:"omitempty,max=500"`
	Variants     *[]ProductVariantRequest `json:"variants,omitempty"`
}

// ProductResponse represents the response for product operations
type ProductResponse struct {
	ID           uint                     `json:"id"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	VendorID     uint                     `json:"vendor_id"`
	CategoryID   uint                     `json:"category_id"`
	CategoryName string                   `json:"category_name"`
	IsActive     bool                     `json:"is_active"`
	IsFeatured   bool                     `json:"is_featured"`
	ImageURLs    []string                 `json:"image_urls"`
	ThumbnailURL string                   `json:"thumbnail_url,omitempty"`
	Variants     []ProductVariantResponse `json:"variants,omitempty"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
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
	CategoryID        uint   `query:"category_id" validate:"omitempty"`
	CategoryName      string `query:"category_name" validate:"omitempty,max=100"`
	ProductID         uint   `query:"product_id" validate:"omitempty"`
	IsActive          *bool  `query:"is_active" validate:"omitempty"`
	IsFeatured        *bool  `query:"is_featured" validate:"omitempty"`
	NameOrDescription string `query:"name_or_description" validate:"omitempty,max=100"`
	Limit             int    `query:"limit" validate:"min=1,max=100"`
	Offset            int    `query:"offset" validate:"min=0"`
	Ordering          string `query:"ordering" validate:"omitempty,oneof=name price created_at updated_at"`
}
