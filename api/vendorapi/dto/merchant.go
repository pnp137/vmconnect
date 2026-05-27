package dto

import "time"

// AddMerchantRequest represents a merchant created by a vendor.
type AddMerchantRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=100"`
	Email        string `json:"email" validate:"required,email"`
	Phone        string `json:"phone" validate:"required,min=10,max=15"`
	Password     string `json:"password" validate:"required,min=6"`
	BusinessName string `json:"business_name" validate:"required,min=2,max=100"`
	ShopName     string `json:"shop_name" validate:"required,min=2,max=100"`
	Address      string `json:"address" validate:"required"`
	Pincode      string `json:"pincode" validate:"required,min=6,max=6"`
}

// AddMerchantRequestDto represents the internal DTO for vendor-created merchants.
type AddMerchantRequestDto struct {
	Name         string
	Email        string
	Phone        string
	Password     string
	BusinessName string
	ShopName     string
	Address      string
	Pincode      string
}

// UpdateMerchantRequest represents merchant fields a vendor can update.
type UpdateMerchantRequest struct {
	Name         *string `json:"name,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Password     *string `json:"password,omitempty"`
	BusinessName *string `json:"business_name,omitempty"`
	ShopName     *string `json:"shop_name,omitempty"`
	Address      *string `json:"address,omitempty"`
	Pincode      *string `json:"pincode,omitempty"`
}

// UpdateMerchantRequestDto represents the internal DTO for vendor merchant updates.
type UpdateMerchantRequestDto struct {
	Name         *string
	Email        *string
	Phone        *string
	Password     *string
	BusinessName *string
	ShopName     *string
	Address      *string
	Pincode      *string
}

// MerchantInfo represents merchant information for vendor
type MerchantInfo struct {
	ID           uint      `json:"id"`
	OwnerName    string    `json:"owner_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	BusinessName string    `json:"business_name"`
	ShopName     string    `json:"shop_name"`
	Address      string    `json:"address"`
	Pincode      string    `json:"pincode"`
	CreatedAt    time.Time `json:"created_at"`
}

// MerchantQueryParam represents query parameters for paginated merchant listing.
type MerchantQueryParam struct {
	BusinessName string `query:"business_name" validate:"omitempty,max=100"`
	ShopName     string `query:"shop_name" validate:"omitempty,max=100"`
	Pincode      string `query:"pincode" validate:"omitempty,max=10"`
	OwnerName    string `query:"owner_name" validate:"omitempty,max=100"`
	MerchantID   uint   `query:"merchant_id" validate:"omitempty"`
	Limit        int    `query:"limit" validate:"min=1,max=100"`
	Offset       int    `query:"offset" validate:"min=0"`
	Ordering     string `query:"ordering" validate:"omitempty,oneof=owner_name business_name shop_name created_at"`
}

// GetMerchantsResponse represents the response for getting vendor's merchants
type GetMerchantsResponse struct {
	Merchants []MerchantInfo `json:"merchants"`
	Count     int            `json:"count"`
	Limit     int            `json:"limit"`
	Offset    int            `json:"offset"`
	Total     int            `json:"total"`
}

// MerchantResponse represents create/update/delete merchant responses.
type MerchantResponse struct {
	Merchant MerchantInfo `json:"merchant"`
	Message  string       `json:"message"`
}

// DeleteMerchantResponse represents the response for deleting a merchant.
type DeleteMerchantResponse struct {
	Message string `json:"message"`
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

// ToAddMerchantRequestDto converts AddMerchantRequest to AddMerchantRequestDto.
func ToAddMerchantRequestDto(req *AddMerchantRequest) *AddMerchantRequestDto {
	return &AddMerchantRequestDto{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password,
		BusinessName: req.BusinessName,
		ShopName:     req.ShopName,
		Address:      req.Address,
		Pincode:      req.Pincode,
	}
}

// ToUpdateMerchantRequestDto converts UpdateMerchantRequest to UpdateMerchantRequestDto.
func ToUpdateMerchantRequestDto(req *UpdateMerchantRequest) *UpdateMerchantRequestDto {
	return &UpdateMerchantRequestDto{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Password:     req.Password,
		BusinessName: req.BusinessName,
		ShopName:     req.ShopName,
		Address:      req.Address,
		Pincode:      req.Pincode,
	}
}
