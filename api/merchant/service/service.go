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

// RegisterMerchant handles merchant registration logic.
func (s *MerchantServiceImpl) RegisterMerchant(ctx context.Context, registerDto *dto.MerchantRegisterRequestDto) (*dto.MerchantResponse, *response.ErrorDetails) {
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

	hashedPassword, err := auth.HashPassword(registerDto.Password)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process password",
			Error:   err,
		}
	}

	merchantRole := models.ROLE_MERCHANT
	user := &models.User{
		Name:         registerDto.Name,
		Email:        registerDto.Email,
		Phone:        registerDto.Phone,
		PasswordHash: hashedPassword,
		RoleID:       &merchantRole,
		IsActive:     true,
	}

	createdUser, err := s.merchantRepository.CreateUser(user)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err,
		}
	}

	merchant := &models.Merchant{
		OwnerID:      createdUser.ID,
		BusinessName: registerDto.BusinessName,
		ShopName:     registerDto.ShopName,
		Address:      registerDto.Address,
		Pincode:      registerDto.Pincode,
	}

	createdMerchant, err := s.merchantRepository.CreateMerchant(merchant)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create merchant profile",
			Error:   err,
		}
	}

	return &dto.MerchantResponse{
		User: dto.MerchantUserInfo{
			ID:           createdUser.ID,
			Name:         createdUser.Name,
			Email:        createdUser.Email,
			Phone:        createdUser.Phone,
			BusinessName: createdMerchant.BusinessName,
			ShopName:     createdMerchant.ShopName,
			Address:      createdMerchant.Address,
			Pincode:      createdMerchant.Pincode,
			Role:         models.ROLE_MERCHANT.String(),
		},
		Message: "Merchant registered successfully",
	}, nil
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
		ID:           merchant.Owner.ID,
		Name:         merchant.Owner.Name,
		Email:        merchant.Owner.Email,
		Phone:        merchant.Owner.Phone,
		BusinessName: merchant.BusinessName,
		ShopName:     merchant.ShopName,
		Address:      merchant.Address,
		Pincode:      merchant.Pincode,
		Role:         models.ROLE_MERCHANT.String(),
	}

	return merchantInfo, nil
}
