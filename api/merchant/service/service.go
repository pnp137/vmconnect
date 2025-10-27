package service

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/repository"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils/auth"
)

// MerchantService defines the interface for merchant service
type MerchantService interface {
	RegisterMerchant(ctx context.Context, registerDto *dto.MerchantRegisterRequestDto) (*dto.MerchantResponse, *response.ErrorDetails)
	GetMerchantInfo(ctx context.Context, merchantID uint) (*dto.MerchantUserInfo, *response.ErrorDetails)
}

type MerchantServiceImpl struct {
	merchantRepository repository.MerchantRepository
}

// NewMerchantService creates a new instance of merchant service
func NewMerchantService(repo repository.MerchantRepository) MerchantService {
	return &MerchantServiceImpl{
		merchantRepository: repo,
	}
}

// RegisterMerchant handles merchant registration logic
func (s *MerchantServiceImpl) RegisterMerchant(ctx context.Context, registerDto *dto.MerchantRegisterRequestDto) (*dto.MerchantResponse, *response.ErrorDetails) {
	// Check if user already exists
	_, err := s.merchantRepository.GetUserByEmailOrPhone(registerDto.Email, registerDto.Phone)
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

	// Validate vendor code if provided
	if registerDto.VendorCode != "" {
		_, err := s.merchantRepository.GetAndValidateVendorByCode(registerDto.VendorCode)
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
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(registerDto.Password)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process password",
			Error:   err,
		}
	}

	// Set merchant role
	merchantRole := models.ROLE_MERCHANT

	// Create user
	user := &models.User{
		Name:         registerDto.Name,
		Email:        registerDto.Email,
		Phone:        registerDto.Phone,
		PasswordHash: hashedPassword,
		RoleID:       &merchantRole,
	}

	createdUser, err := s.merchantRepository.CreateUser(user)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err,
		}
	}

	// Create merchant profile
	merchant := &models.Merchant{
		UserID:     createdUser.ID,
		VendorCode: registerDto.VendorCode,
		ShopName:   registerDto.ShopName,
		Address:    registerDto.Address,
		Pincode:    registerDto.Pincode,
	}

	createdMerchant, err := s.merchantRepository.CreateMerchant(merchant)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create merchant profile",
			Error:   err,
		}
	}

	// Create vendor-merchant mapping if vendor code is provided
	if registerDto.VendorCode != "" {
		// Get vendor by code (already validated above)
		vendor, err := s.merchantRepository.GetAndValidateVendorByCode(registerDto.VendorCode)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get vendor information",
				Error:   err,
			}
		}

		// Create vendor-merchant mapping
		vendorMerchantMapping := &models.VendorMerchantMapping{
			VendorID:   vendor.ID,
			MerchantID: createdMerchant.ID,
			Status:     "active",
		}

		_, err = s.merchantRepository.CreateVendorMerchantMapping(vendorMerchantMapping)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create vendor-merchant relationship",
				Error:   err,
			}
		}
	}

	// Prepare response
	response := &dto.MerchantResponse{
		User: dto.MerchantUserInfo{
			ID:       createdUser.ID,
			Name:     createdUser.Name,
			Email:    createdUser.Email,
			Phone:    createdUser.Phone,
			ShopName: createdMerchant.ShopName,
			Address:  createdMerchant.Address,
			Pincode:  createdMerchant.Pincode,
			Role:     "merchant",
		},
		Message: "Merchant registered successfully",
	}

	return response, nil
}

// GetMerchantInfo gets merchant information by ID
func (s *MerchantServiceImpl) GetMerchantInfo(ctx context.Context, merchantID uint) (*dto.MerchantUserInfo, *response.ErrorDetails) {
	// Get merchant with user info
	merchant, err := s.merchantRepository.GetMerchantByID(merchantID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusNotFound,
				Message: "Merchant not found",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get merchant information",
			Error:   err,
		}
	}

	// Convert to DTO
	merchantInfo := &dto.MerchantUserInfo{
		ID:       merchant.User.ID,
		Name:     merchant.User.Name,
		Email:    merchant.User.Email,
		Phone:    merchant.User.Phone,
		ShopName: merchant.ShopName,
		Address:  merchant.Address,
		Pincode:  merchant.Pincode,
		Role:     "merchant",
	}

	return merchantInfo, nil
}
