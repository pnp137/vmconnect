package service

import (
	"context"
	"net/http"
	"strconv"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/repository"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
)

// MerchantVendorService defines the interface for merchant-vendor operations
type MerchantVendorService interface {
	GetMerchantVendors(ctx context.Context, merchantID uint) (*dto.GetVendorsResponse, *response.ErrorDetails)
	AddVendorToMerchant(ctx context.Context, merchantID uint, req *dto.AddVendorRequest) (*dto.AddVendorResponse, *response.ErrorDetails)
	GetVendorProducts(ctx context.Context, merchantID, vendorID uint) (*dto.GetProductsResponse, *response.ErrorDetails)
	GetVendorCategories(ctx context.Context, merchantID, vendorID uint) (*dto.GetCategoriesResponse, *response.ErrorDetails)
	GetVendorProductsByCategory(ctx context.Context, merchantID, vendorID uint, queryParams *dto.ProductQueryParam) (*dto.GetPaginatedProductsResponse, *response.ErrorDetails)
}

type MerchantVendorServiceImpl struct {
	repository repository.MerchantRepository
}

// NewMerchantVendorService creates a new instance of merchant vendor service
func NewMerchantVendorService(repo repository.MerchantRepository) MerchantVendorService {
	return &MerchantVendorServiceImpl{
		repository: repo,
	}
}

// GetMerchantVendors gets all vendors linked to a merchant
func (s *MerchantVendorServiceImpl) GetMerchantVendors(ctx context.Context, merchantID uint) (*dto.GetVendorsResponse, *response.ErrorDetails) {
	vendors, err := s.repository.GetMerchantVendors(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant vendors",
			Error:   err,
		}
	}

	// Convert to DTO
	vendorInfos := make([]dto.VendorInfo, len(vendors))
	for i, vendor := range vendors {
		vendorInfos[i] = dto.VendorInfo{
			ID:          vendor.ID,
			VendorCode:  vendor.VendorCode,
			CompanyName: vendor.CompanyName,
			GSTNumber:   vendor.GSTNumber,
			Address:     vendor.Address,
			LogoURL:     vendor.LogoURL,
		}
	}

	return &dto.GetVendorsResponse{
		Vendors: vendorInfos,
		Count:   len(vendorInfos),
	}, nil
}

// AddVendorToMerchant adds a vendor to merchant using vendor code
func (s *MerchantVendorServiceImpl) AddVendorToMerchant(ctx context.Context, merchantID uint, req *dto.AddVendorRequest) (*dto.AddVendorResponse, *response.ErrorDetails) {
	// Get vendor by code
	vendor, err := s.repository.GetAndValidateVendorByCode(req.VendorCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusBadRequest,
				Message: "Invalid vendor code provided",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Error validating vendor code",
			Error:   err,
		}
	}

	// Check if mapping already exists
	existingVendors, err := s.repository.GetMerchantVendors(merchantID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Error checking existing vendors",
			Error:   err,
		}
	}

	// Check if vendor is already linked
	for _, existingVendor := range existingVendors {
		if existingVendor.ID == vendor.ID {
			return nil, &response.ErrorDetails{
				Code:    http.StatusConflict,
				Message: "Vendor is already linked to this merchant",
				Error:   nil,
			}
		}
	}

	// Create vendor-merchant mapping
	vendorMerchantMapping := &models.VendorMerchantMapping{
		VendorID:   vendor.ID,
		MerchantID: merchantID,
		Status:     "active",
	}

	_, err = s.repository.CreateVendorMerchantMapping(vendorMerchantMapping)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create vendor-merchant relationship",
			Error:   err,
		}
	}

	// Prepare response
	vendorInfo := dto.VendorInfo{
		ID:          vendor.ID,
		VendorCode:  vendor.VendorCode,
		CompanyName: vendor.CompanyName,
		GSTNumber:   vendor.GSTNumber,
		Address:     vendor.Address,
		LogoURL:     vendor.LogoURL,
	}

	return &dto.AddVendorResponse{
		Vendor:  vendorInfo,
		Message: "Vendor added successfully",
	}, nil
}

// GetVendorProducts gets products for a specific vendor-merchant relationship
func (s *MerchantVendorServiceImpl) GetVendorProducts(ctx context.Context, merchantID, vendorID uint) (*dto.GetProductsResponse, *response.ErrorDetails) {
	// Get products using optimized repository method
	products, err := s.repository.GetVendorProductsForMerchant(merchantID, vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor products",
			Error:   err,
		}
	}

	// Group products by category
	categories := make(map[string][]dto.ProductInfo)
	totalCount := 0

	for _, product := range products {
		categoryName := product.Category.Name
		variant := getPrimaryVariant(&product)
		sku := 0
		unit := ""
		price := 0.0
		stock := 0
		if variant != nil {
			sku, _ = strconv.Atoi(variant.Name)
			unit = variant.Unit
			price = variant.Price
			stock = variant.Stock
		}
		productInfo := dto.ProductInfo{
			ID:           product.ID,
			Name:         product.Name,
			Description:  product.Description,
			Category:     categoryName,
			SKU:          sku,
			Unit:         unit,
			Price:        price,
			Stock:        stock,
			ImageURLs:    []string(product.ImageURLs),
			ThumbnailURL: "",
			IsActive:     product.IsActive,
			IsFeatured:   product.IsFeatured,
		}

		categories[categoryName] = append(categories[categoryName], productInfo)
		totalCount++
	}

	return &dto.GetProductsResponse{
		Categories: categories,
		Count:      totalCount,
	}, nil
}

// GetVendorCategories gets unique categories for a specific vendor-merchant relationship
func (s *MerchantVendorServiceImpl) GetVendorCategories(ctx context.Context, merchantID, vendorID uint) (*dto.GetCategoriesResponse, *response.ErrorDetails) {
	// Get categories using repository method
	categories, err := s.repository.GetVendorCategoriesForMerchant(merchantID, vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor categories",
			Error:   err,
		}
	}

	// Convert to DTO
	categoryInfos := make([]dto.CategoryInfo, len(categories))
	for i, category := range categories {
		categoryInfos[i] = dto.CategoryInfo{
			ID:   category.ID,
			Name: category.Name,
		}
	}

	return &dto.GetCategoriesResponse{
		Categories: categoryInfos,
		Count:      len(categoryInfos),
	}, nil
}

// GetVendorProductsByCategory gets paginated products for a specific category
func (s *MerchantVendorServiceImpl) GetVendorProductsByCategory(ctx context.Context, merchantID, vendorID uint, queryParams *dto.ProductQueryParam) (*dto.GetPaginatedProductsResponse, *response.ErrorDetails) {
	// Get products using repository method with pagination
	products, totalCount, err := s.repository.GetVendorProductsByCategoryForMerchant(merchantID, vendorID, queryParams)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor products by category",
			Error:   err,
		}
	}

	// Convert to DTO
	productInfos := make([]dto.ProductInfo, len(products))
	for i, product := range products {
		categoryName := product.Category.Name
		variant := getPrimaryVariant(&product)
		sku := 0
		unit := ""
		price := 0.0
		stock := 0
		if variant != nil {
			sku, _ = strconv.Atoi(variant.Name)
			unit = variant.Unit
			price = variant.Price
			stock = variant.Stock
		}
		productInfos[i] = dto.ProductInfo{
			ID:           product.ID,
			Name:         product.Name,
			Description:  product.Description,
			Category:     categoryName,
			SKU:          sku,
			Unit:         unit,
			Price:        price,
			Stock:        stock,
			ImageURLs:    []string(product.ImageURLs),
			ThumbnailURL: "",
			IsActive:     product.IsActive,
			IsFeatured:   product.IsFeatured,
		}
	}

	return &dto.GetPaginatedProductsResponse{
		Products: productInfos,
		Count:    len(productInfos),
		Limit:    queryParams.Limit,
		Offset:   queryParams.Offset,
		Total:    totalCount,
	}, nil
}

func getPrimaryVariant(product *models.Product) *models.ProductVariant {
	if product == nil || len(product.Variants) == 0 {
		return nil
	}
	return &product.Variants[0]
}
