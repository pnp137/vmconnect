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

	// Get merchant role
	merchantRole, err := s.merchantRepository.GetRoleByName("merchant")
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Merchant role not found",
			Error:   err,
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

	// Create user
	user := &models.User{
		Name:         registerDto.Name,
		Email:        registerDto.Email,
		Phone:        registerDto.Phone,
		PasswordHash: hashedPassword,
		RoleID:       merchantRole.ID,
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
		UserID:   createdUser.ID,
		ShopName: registerDto.ShopName,
		Address:  registerDto.Address,
		Pincode:  registerDto.Pincode,
	}

	createdMerchant, err := s.merchantRepository.CreateMerchant(merchant)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create merchant profile",
			Error:   err,
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
