package service

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	validator "linksupply.io/vmconnect/api/vendorapi/validator"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils/auth"
)

type VendorMerchantService interface {
	GetMerchants(ctx context.Context, vendorID uint, queryParams *dto.MerchantQueryParam) (*dto.GetMerchantsResponse, *response.ErrorDetails)
	AddMerchant(ctx context.Context, vendorID uint, req *dto.AddMerchantRequestDto) (*dto.MerchantResponse, *response.ErrorDetails)
	UpdateMerchant(ctx context.Context, vendorID, merchantID uint, req *dto.UpdateMerchantRequestDto) (*dto.MerchantResponse, *response.ErrorDetails)
	DeleteMerchant(ctx context.Context, vendorID, merchantID uint) (*dto.DeleteMerchantResponse, *response.ErrorDetails)
	GetProductVisibility(ctx context.Context, vendorID, merchantID uint) (*dto.GetProductVisibilityResponse, *response.ErrorDetails)
	UpdateProductVisibility(ctx context.Context, vendorID, merchantID uint, req []dto.UpdateProductVisibilityRequest) (*dto.UpdateProductVisibilityResponse, *response.ErrorDetails)
}

type VendorMerchantServiceImpl struct {
	repository     repository.VendorMerchantRepository
	listRepository repository.VendorMerchantListRepository
	validator      validator.VendorValidator
}

func NewVendorMerchantService(repo repository.VendorMerchantRepository, listRepo repository.VendorMerchantListRepository) VendorMerchantService {
	return &VendorMerchantServiceImpl{
		repository:     repo,
		listRepository: listRepo,
		validator:      validator.NewVendorValidator(),
	}
}

// GetMerchants gets merchants linked to a vendor with pagination
func (s *VendorMerchantServiceImpl) GetMerchants(ctx context.Context, vendorID uint, queryParams *dto.MerchantQueryParam) (*dto.GetMerchantsResponse, *response.ErrorDetails) {
	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}
	if queryParams.Ordering == "" {
		queryParams.Ordering = "created_at"
	}

	if err := s.validator.ValidateMerchantQueryParam(queryParams); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Error:   err,
		}
	}

	merchants, totalCount, err := s.listRepository.GetVendorMerchants(vendorID, queryParams)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchants",
			Error:   err,
		}
	}

	merchantResponses := make([]dto.MerchantInfo, len(merchants))
	for i, merchant := range merchants {
		merchantResponses[i] = toMerchantInfo(&merchant)
	}

	return &dto.GetMerchantsResponse{
		Merchants: merchantResponses,
		Count:     len(merchantResponses),
		Limit:     queryParams.Limit,
		Offset:    queryParams.Offset,
		Total:     totalCount,
	}, nil
}

// AddMerchant creates a merchant owned by the given vendor.
func (s *VendorMerchantServiceImpl) AddMerchant(ctx context.Context, vendorID uint, req *dto.AddMerchantRequestDto) (*dto.MerchantResponse, *response.ErrorDetails) {
	if _, err := s.repository.GetVendorByID(vendorID); err != nil {
		return nil, vendorNotFoundOrError(err)
	}

	_, err := s.repository.GetUserByEmailOrPhone(req.Email, req.Phone)
	if err == nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusConflict,
			Message: "User already exists with this email or phone",
			Error:   nil,
		}
	}
	if err != gorm.ErrRecordNotFound {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Error checking existing user",
			Error:   err,
		}
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process password",
			Error:   err,
		}
	}

	merchantRole := models.ROLE_MERCHANT
	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		RoleID:       &merchantRole,
		IsActive:     true,
	}

	merchant := &models.Merchant{
		BusinessName: req.BusinessName,
		ShopName:     req.ShopName,
		Address:      req.Address,
		Pincode:      req.Pincode,
	}

	mapping := &models.VendorMerchantMapping{
		VendorID: vendorID,
		Status:   "active",
	}

	createdMerchant, err := s.repository.CreateMerchantForVendor(user, merchant, mapping)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create merchant",
			Error:   err,
		}
	}

	return &dto.MerchantResponse{
		Merchant: toMerchantInfo(createdMerchant),
		Message:  "Merchant added successfully",
	}, nil
}

// UpdateMerchant updates a merchant owned by the given vendor.
func (s *VendorMerchantServiceImpl) UpdateMerchant(ctx context.Context, vendorID, merchantID uint, req *dto.UpdateMerchantRequestDto) (*dto.MerchantResponse, *response.ErrorDetails) {
	merchant, err := s.repository.GetVendorMerchant(vendorID, merchantID)
	if err != nil {
		return nil, merchantNotFoundOrError(err)
	}

	email := merchant.Owner.Email
	phone := merchant.Owner.Phone
	if req.Email != nil {
		email = *req.Email
	}
	if req.Phone != nil {
		phone = *req.Phone
	}

	if req.Email != nil || req.Phone != nil {
		_, err := s.repository.GetUserByEmailOrPhoneExcludingUser(email, phone, merchant.OwnerID)
		if err == nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusConflict,
				Message: "User already exists with this email or phone",
				Error:   nil,
			}
		}
		if err != gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Error checking existing user",
				Error:   err,
			}
		}
	}

	userUpdates := map[string]interface{}{}
	merchantUpdates := map[string]interface{}{}

	if req.Name != nil {
		userUpdates["name"] = *req.Name
	}
	if req.Email != nil {
		userUpdates["email"] = *req.Email
	}
	if req.Phone != nil {
		userUpdates["phone"] = *req.Phone
	}
	if req.Password != nil {
		hashedPassword, err := auth.HashPassword(*req.Password)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to process password",
				Error:   err,
			}
		}
		userUpdates["password_hash"] = hashedPassword
	}
	if req.BusinessName != nil {
		merchantUpdates["business_name"] = *req.BusinessName
	}
	if req.ShopName != nil {
		merchantUpdates["shop_name"] = *req.ShopName
	}
	if req.Address != nil {
		merchantUpdates["address"] = *req.Address
	}
	if req.Pincode != nil {
		merchantUpdates["pincode"] = *req.Pincode
	}

	updatedMerchant, err := s.repository.UpdateMerchantForVendor(merchant, userUpdates, merchantUpdates)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update merchant",
			Error:   err,
		}
	}

	return &dto.MerchantResponse{
		Merchant: toMerchantInfo(updatedMerchant),
		Message:  "Merchant updated successfully",
	}, nil
}

// DeleteMerchant deletes a merchant owned by the given vendor.
func (s *VendorMerchantServiceImpl) DeleteMerchant(ctx context.Context, vendorID, merchantID uint) (*dto.DeleteMerchantResponse, *response.ErrorDetails) {
	err := s.repository.DeleteMerchantForVendor(vendorID, merchantID)
	if err != nil {
		return nil, merchantNotFoundOrError(err)
	}

	return &dto.DeleteMerchantResponse{
		Message: "Merchant deleted successfully",
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
			Price:        getProductPrimaryPrice(&visibility.Product),
			Visible:      !visibility.IsDeleted, // visible = !is_deleted
		}
	}

	return &dto.GetProductVisibilityResponse{
		Products: visibilityResponses,
		Count:    len(visibilityResponses),
	}, nil
}

func getProductPrimaryPrice(product *models.Product) float64 {
	if product == nil || len(product.Variants) == 0 {
		return 0
	}
	return product.Variants[0].Price
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

func toMerchantInfo(merchant *models.Merchant) dto.MerchantInfo {
	return dto.MerchantInfo{
		ID:           merchant.ID,
		OwnerName:    merchant.Owner.Name,
		Email:        merchant.Owner.Email,
		Phone:        merchant.Owner.Phone,
		BusinessName: merchant.BusinessName,
		ShopName:     merchant.ShopName,
		Address:      merchant.Address,
		Pincode:      merchant.Pincode,
		CreatedAt:    merchant.CreatedAt,
	}
}

func vendorNotFoundOrError(err error) *response.ErrorDetails {
	if err == gorm.ErrRecordNotFound {
		return &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Vendor not found",
			Error:   err,
		}
	}

	return &response.ErrorDetails{
		Code:    http.StatusInternalServerError,
		Message: "Failed to get vendor",
		Error:   err,
	}
}

func merchantNotFoundOrError(err error) *response.ErrorDetails {
	if err == gorm.ErrRecordNotFound {
		return &response.ErrorDetails{
			Code:    http.StatusNotFound,
			Message: "Merchant not found for this vendor",
			Error:   err,
		}
	}

	return &response.ErrorDetails{
		Code:    http.StatusInternalServerError,
		Message: "Failed to get merchant",
		Error:   err,
	}
}
