package service

import (
	"context"
	"net/http"

	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
)

type VendorMerchantService interface {
	GetMerchants(ctx context.Context, vendorID uint) (*dto.GetMerchantsResponse, *response.ErrorDetails)
	GetProductVisibility(ctx context.Context, vendorID, merchantID uint) (*dto.GetProductVisibilityResponse, *response.ErrorDetails)
	UpdateProductVisibility(ctx context.Context, vendorID, merchantID uint, req []dto.UpdateProductVisibilityRequest) (*dto.UpdateProductVisibilityResponse, *response.ErrorDetails)
}

type VendorMerchantServiceImpl struct {
	repository repository.VendorMerchantRepository
	validator  validator.VendorValidator
}

func NewVendorMerchantService(repo repository.VendorMerchantRepository) VendorMerchantService {
	return &VendorMerchantServiceImpl{
		repository: repo,
		validator:  validator.NewVendorValidator(),
	}
}

// GetMerchants gets all merchants linked to a vendor
func (s *VendorMerchantServiceImpl) GetMerchants(ctx context.Context, vendorID uint) (*dto.GetMerchantsResponse, *response.ErrorDetails) {
	// Get merchants from repository
	merchants, err := s.repository.GetVendorMerchants(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchants",
			Error:   err,
		}
	}

	// Convert to response DTOs
	merchantResponses := make([]dto.MerchantInfo, len(merchants))
	for i, merchant := range merchants {
		merchantResponses[i] = dto.MerchantInfo{
			ID:        merchant.ID,
			Name:      merchant.User.Name,
			Email:     merchant.User.Email,
			Phone:     merchant.User.Phone,
			ShopName:  merchant.ShopName,
			Address:   merchant.Address,
			Pincode:   merchant.Pincode,
			CreatedAt: merchant.CreatedAt,
		}
	}

	return &dto.GetMerchantsResponse{
		Merchants: merchantResponses,
		Count:     len(merchantResponses),
	}, nil
}

// GetProductVisibility gets product visibility settings for a specific merchant
func (s *VendorMerchantServiceImpl) GetProductVisibility(ctx context.Context, vendorID, merchantID uint) (*dto.GetProductVisibilityResponse, *response.ErrorDetails) {
	// Get product visibilities from repository
	visibilities, err := s.repository.GetProductVisibilityForMerchant(vendorID, merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get product visibility",
			Error:   err,
		}
	}

	// Convert to response DTOs
	visibilityResponses := make([]dto.ProductVisibilityInfo, len(visibilities))
	for i, visibility := range visibilities {
		categoryName := ""
		if visibility.Product.CategoryID > 0 && visibility.Product.Category.ID > 0 {
			categoryName = visibility.Product.Category.Name
		}

		visibilityResponses[i] = dto.ProductVisibilityInfo{
			ProductID:    visibility.ProductID,
			ProductName:  visibility.Product.Name,
			CategoryID:   visibility.Product.CategoryID,
			CategoryName: categoryName,
			Price:        visibility.Product.Price,
			Visible:      !visibility.IsDeleted, // visible = !is_deleted
		}
	}

	return &dto.GetProductVisibilityResponse{
		Products: visibilityResponses,
		Count:    len(visibilityResponses),
	}, nil
}

// UpdateProductVisibility updates product visibility settings for a merchant
func (s *VendorMerchantServiceImpl) UpdateProductVisibility(ctx context.Context, vendorID, merchantID uint, req []dto.UpdateProductVisibilityRequest) (*dto.UpdateProductVisibilityResponse, *response.ErrorDetails) {
	// Validate request
	if len(req) == 0 {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: "At least one product visibility update is required",
			Error:   nil,
		}
	}

	// Update product visibilities in repository
	err := s.repository.UpdateProductVisibilityForMerchant(vendorID, merchantID, req)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update product visibility",
			Error:   err,
		}
	}

	return &dto.UpdateProductVisibilityResponse{
		UpdatedCount: len(req),
		Message:      "Product visibility updated successfully",
	}, nil
}
