package service

import (
	"context"
	"net/http"
	"strings"

	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/database/models"
)

type VendorProductService interface {
	AddProduct(ctx context.Context, vendorID uint, req *dto.AddProductRequest) (*dto.AddProductResponse, *response.ErrorDetails)
	UpdateProduct(ctx context.Context, vendorID, productID uint, req *dto.UpdateProductRequest) (*dto.UpdateProductResponse, *response.ErrorDetails)
	DeleteProduct(ctx context.Context, vendorID, productID uint) (*dto.DeleteProductResponse, *response.ErrorDetails)
	GetProduct(ctx context.Context, vendorID, productID uint) (*dto.ProductResponse, *response.ErrorDetails)
	GetProducts(ctx context.Context, vendorID uint, queryParams *dto.ProductQueryParam) (*dto.GetProductsResponse, *response.ErrorDetails)
	GetCategories(ctx context.Context, vendorID uint) (*dto.GetCategoriesResponse, *response.ErrorDetails)
}

type VendorProductServiceImpl struct {
	repository repository.VendorProductRepository
	validator  validator.VendorValidator
}

func NewVendorProductService(repo repository.VendorProductRepository) VendorProductService {
	return &VendorProductServiceImpl{
		repository: repo,
		validator:  validator.NewVendorValidator(),
	}
}

// AddProduct adds a new product for a vendor
func (s *VendorProductServiceImpl) AddProduct(ctx context.Context, vendorID uint, req *dto.AddProductRequest) (*dto.AddProductResponse, *response.ErrorDetails) {
	// Validate request
	if err := s.validator.ValidateAddProductRequest(req); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	category, err := s.repository.GetOrCreateCategory(vendorID, req.CategoryName)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get or create category",
			Error:   err,
		}
	}

	// Create product model
	product := &models.Product{
		Name:         req.Name,
		Description:  req.Description,
		CategoryID:   category.ID,
		IsActive:     req.IsActive,
		IsFeatured:   req.IsFeatured,
		ImageURLs:    cleanProductImages(req.ImageURLs),
		ThumbnailURL: req.ThumbnailURL,
	}

	// Create product in repository
	createdProduct, err := s.repository.CreateProduct(vendorID, product)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create product",
			Error:   err,
		}
	}

	var variants []models.ProductVariant
	for _, v := range req.Variants {
		active := true
		if v.IsActive != nil {
			active = *v.IsActive
		}
		variants = append(variants, models.ProductVariant{
			Name:     v.Name,
			Unit:     v.Unit,
			Quantity: v.Quantity,
			Price:    v.Price,
			MRP:      v.MRP,
			MOQ:      v.MOQ,
			Stock:    v.Stock,
			IsActive: active,
		})
	}
	if err := s.repository.CreateProductVariants(vendorID, createdProduct.ID, variants); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create product variants",
			Error:   err,
		}
	}
	createdProduct, err = s.repository.GetProductByID(vendorID, createdProduct.ID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to load created product",
			Error:   err,
		}
	}

	// Convert to response DTO
	productResponse := s.convertToProductResponse(createdProduct)

	return &dto.AddProductResponse{
		Product: productResponse,
		Message: "Product created successfully",
	}, nil
}

// UpdateProduct updates a product for a vendor
func (s *VendorProductServiceImpl) UpdateProduct(ctx context.Context, vendorID, productID uint, req *dto.UpdateProductRequest) (*dto.UpdateProductResponse, *response.ErrorDetails) {
	// Validate request
	if err := s.validator.ValidateUpdateProductRequest(req); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Convert request to updates map
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.CategoryName != nil {
		category, err := s.repository.GetOrCreateCategory(vendorID, *req.CategoryName)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get or create category",
				Error:   err,
			}
		}
		updates["category_id"] = category.ID
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}
	if req.ImageURLs != nil {
		updates["image_urls"] = cleanProductImages(*req.ImageURLs)
	}
	if req.ThumbnailURL != nil {
		updates["thumbnail_url"] = *req.ThumbnailURL
	}

	// Update product in repository
	updatedProduct, err := s.repository.UpdateProduct(vendorID, productID, updates)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update product",
			Error:   err,
		}
	}

	if req.Variants != nil {
		var variants []models.ProductVariant
		for _, v := range *req.Variants {
			active := true
			if v.IsActive != nil {
				active = *v.IsActive
			}
			var id uint
			if v.ID != nil {
				id = *v.ID
			}
			variants = append(variants, models.ProductVariant{
				ID:        id,
				ProductID: productID,
				Name:      v.Name,
				Unit:      v.Unit,
				Quantity:  v.Quantity,
				Price:     v.Price,
				MRP:       v.MRP,
				MOQ:       v.MOQ,
				Stock:     v.Stock,
				IsActive:  active,
			})
		}

		if err := s.repository.SyncProductVariants(vendorID, productID, variants); err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to sync product variants",
				Error:   err,
			}
		}
		updatedProduct, err = s.repository.GetProductByID(vendorID, productID)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to load updated product",
				Error:   err,
			}
		}
	}

	// Convert to response DTO
	productResponse := s.convertToProductResponse(updatedProduct)

	return &dto.UpdateProductResponse{
		Product: productResponse,
		Message: "Product updated successfully",
	}, nil
}

// DeleteProduct deletes a product for a vendor
func (s *VendorProductServiceImpl) DeleteProduct(ctx context.Context, vendorID, productID uint) (*dto.DeleteProductResponse, *response.ErrorDetails) {
	// Delete product in repository
	err := s.repository.DeleteProduct(vendorID, productID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to delete product",
			Error:   err,
		}
	}

	return &dto.DeleteProductResponse{
		ProductID: productID,
		Message:   "Product deleted successfully",
	}, nil
}

// GetProduct gets a product by ID for a vendor
func (s *VendorProductServiceImpl) GetProduct(ctx context.Context, vendorID, productID uint) (*dto.ProductResponse, *response.ErrorDetails) {
	// Get product from repository
	product, err := s.repository.GetProductByID(vendorID, productID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Product not found",
			Error:   err,
		}
	}

	// Convert to response DTO
	productResponse := s.convertToProductResponse(product)

	return &productResponse, nil
}

// GetProducts gets products for a vendor with filtering and pagination
func (s *VendorProductServiceImpl) GetProducts(ctx context.Context, vendorID uint, queryParams *dto.ProductQueryParam) (*dto.GetProductsResponse, *response.ErrorDetails) {
	// Validate query parameters
	if err := s.validator.ValidateProductQueryParam(queryParams); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	// Set default values
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Ordering == "" {
		queryParams.Ordering = "created_at"
	}

	// Get products from repository
	products, totalCount, err := s.repository.GetVendorProducts(vendorID, queryParams)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get products",
			Error:   err,
		}
	}

	// Convert to response DTOs
	productResponses := make([]dto.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = s.convertToProductResponse(&product)
	}

	return &dto.GetProductsResponse{
		Products: productResponses,
		Count:    len(productResponses),
		Limit:    queryParams.Limit,
		Offset:   queryParams.Offset,
		Total:    totalCount,
	}, nil
}

// GetCategories gets categories for a vendor
func (s *VendorProductServiceImpl) GetCategories(ctx context.Context, vendorID uint) (*dto.GetCategoriesResponse, *response.ErrorDetails) {
	// Get categories from repository
	categories, err := s.repository.GetVendorCategories(vendorID)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get categories",
			Error:   err,
		}
	}

	// Convert to response DTOs
	categoryResponses := make([]dto.CategoryInfo, len(categories))
	for i, category := range categories {
		categoryResponses[i] = dto.CategoryInfo{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt,
		}
	}

	return &dto.GetCategoriesResponse{
		Categories: categoryResponses,
		Count:      len(categoryResponses),
	}, nil
}

// convertToProductResponse converts a Product model to ProductResponse DTO
func (s *VendorProductServiceImpl) convertToProductResponse(product *models.Product) dto.ProductResponse {
	categoryName := ""
	if product.CategoryID > 0 && product.Category.ID > 0 {
		categoryName = product.Category.Name
	}
	var variantResponses []dto.ProductVariantResponse
	for _, v := range product.Variants {
		variantResponses = append(variantResponses, dto.ProductVariantResponse{
			ID:        v.ID,
			Name:      v.Name,
			Unit:      v.Unit,
			Quantity:  v.Quantity,
			Price:     v.Price,
			MRP:       v.MRP,
			MOQ:       v.MOQ,
			Stock:     v.Stock,
			IsActive:  v.IsActive,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		VendorID:     product.VendorID,
		CategoryID:   product.CategoryID,
		CategoryName: categoryName,
		IsActive:     product.IsActive,
		IsFeatured:   product.IsFeatured,
		ImageURLs:    []string(product.ImageURLs),
		ThumbnailURL: product.ThumbnailURL,
		Variants:     variantResponses,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}

func cleanProductImages(imageURLs []string) models.StringArray {
	clean := make(models.StringArray, 0, len(imageURLs))
	for _, image := range imageURLs {
		trimmed := strings.TrimSpace(image)
		if trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return clean
}
