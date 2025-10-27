package service

import (
	"context"
	"net/http"
	"strconv"

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

	// Get or create category
	var categoryID uint
	if req.CategoryID > 0 {
		// Use provided category ID
		categoryID = req.CategoryID
	} else {
		// Get or create category by name
		category, err := s.repository.GetOrCreateCategory(vendorID, req.CategoryName)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get or create category",
				Error:   err,
			}
		}
		categoryID = category.ID
	}

	// Create product model
	product := &models.Product{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryID:   categoryID,
		CategoryName: req.CategoryName, // Store category name as well
		SKU:          req.SKU,
		IsActive:     req.IsActive,
		IsFeatured:   req.IsFeatured,
		ImageURL:     req.ImageURL,
		Stock:        req.Stock,
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
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.SKU != nil {
		updates["sku"] = *req.SKU
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsFeatured != nil {
		updates["is_featured"] = *req.IsFeatured
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
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

	return dto.ProductResponse{
		ID:           product.ID,
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		CategoryID:   product.CategoryID,
		CategoryName: categoryName,
		SKU:          strconv.Itoa(product.SKU), // Convert int to string for response
		IsActive:     product.IsActive,
		IsFeatured:   product.IsFeatured,
		ImageURL:     product.ImageURL,
		Stock:        product.Stock,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}
