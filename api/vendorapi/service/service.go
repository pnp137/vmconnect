package service

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils/auth"
)

// VendorService defines the interface for vendor service
type VendorService interface {
	RegisterVendor(ctx context.Context, registerDto *dto.VendorRegisterRequestDto) (*dto.VendorResponse, *response.ErrorDetails)
}

type VendorServiceImpl struct {
	repository repository.VendorRepository
}

// NewVendorService creates a new instance of vendor service
func NewVendorService(repo repository.VendorRepository) VendorService {
	return &VendorServiceImpl{
		repository: repo,
	}
}

// RegisterVendor handles vendor registration logic
func (s *VendorServiceImpl) RegisterVendor(ctx context.Context, registerDto *dto.VendorRegisterRequestDto) (*dto.VendorResponse, *response.ErrorDetails) {
	// Check if user already exists
	_, err := s.repository.GetUserByEmailOrPhone(registerDto.Email, registerDto.Phone)
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

	// Get vendor role
	vendorRole, err := s.repository.GetRoleByName("vendor")
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Vendor role not found",
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
		RoleID:       vendorRole.ID,
	}

	createdUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err,
		}
	}

	// Create vendor profile
	vendor := &models.Vendor{
		UserID:      createdUser.ID,
		CompanyName: registerDto.CompanyName,
		GSTNumber:   registerDto.GSTNumber,
		Address:     registerDto.Address,
		LogoURL:     registerDto.LogoURL,
	}

	createdVendor, err := s.repository.CreateVendor(vendor)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create vendor profile",
			Error:   err,
		}
	}

	// Prepare response
	resp := &dto.VendorResponse{
		User: dto.VendorUserInfo{
			ID:          createdUser.ID,
			Name:        createdUser.Name,
			Email:       createdUser.Email,
			Phone:       createdUser.Phone,
			CompanyName: createdVendor.CompanyName,
			GSTNumber:   createdVendor.GSTNumber,
			Address:     createdVendor.Address,
			LogoURL:     createdVendor.LogoURL,
			Role:        "vendor",
		},
		Message: "Vendor registered successfully",
	}

	return resp, nil
}
