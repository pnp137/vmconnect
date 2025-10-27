package dto

// AddVendorRequest represents the request to add a vendor to merchant
type AddVendorRequest struct {
	VendorCode string `json:"vendor_code" validate:"required,min=7,max=7"`
}

// VendorInfo represents vendor information in responses
type VendorInfo struct {
	ID          uint   `json:"id"`
	VendorCode  string `json:"vendor_code"`
	CompanyName string `json:"company_name"`
	GSTNumber   string `json:"gst_number"`
	Address     string `json:"address"`
	LogoURL     string `json:"logo_url"`
}

// AddVendorResponse represents the response when adding a vendor
type AddVendorResponse struct {
	Vendor  VendorInfo `json:"vendor"`
	Message string     `json:"message"`
}

// GetVendorsResponse represents the response for getting merchant's vendors
type GetVendorsResponse struct {
	Vendors []VendorInfo `json:"vendors"`
	Count   int          `json:"count"`
}

// ProductInfo represents product information
type ProductInfo struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	SKU          int     `json:"sku"`
	Unit         string  `json:"unit"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
	ImageURL     string  `json:"image_url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	IsActive     bool    `json:"is_active"`
	IsFeatured   bool    `json:"is_featured"`
	IsWishlisted bool    `json:"is_wishlisted,omitempty"`
}

// GetProductsResponse represents the response for getting vendor products
type GetProductsResponse struct {
	Categories map[string][]ProductInfo `json:"categories"`
	Count      int                      `json:"count"`
}

// CategoryInfo represents category information
type CategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// GetCategoriesResponse represents the response for getting vendor categories
type GetCategoriesResponse struct {
	Categories []CategoryInfo `json:"categories"`
	Count      int            `json:"count"`
}

// ProductQueryParam represents query parameters for paginated product listing
type ProductQueryParam struct {
	CategoryID uint   `query:"category_id" validate:"required"`
	Limit      int    `query:"limit" validate:"min=1,max=100"`
	Offset     int    `query:"offset" validate:"min=0"`
	Ordering   string `query:"ordering" validate:"omitempty,oneof=name price created_at"`
	State      string `query:"state" validate:"omitempty,oneof=active featured"`
}

// GetPaginatedProductsResponse represents the response for paginated product listing
type GetPaginatedProductsResponse struct {
	Products []ProductInfo `json:"products"`
	Count    int           `json:"count"`
	Limit    int           `json:"limit"`
	Offset   int           `json:"offset"`
	Total    int           `json:"total"`
}
