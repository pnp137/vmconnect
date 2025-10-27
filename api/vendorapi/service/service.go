package service

import (
	"context"
	"net/http"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils"
	"linksupply.io/vmconnect/utils/auth"
)

// VendorService defines the interface for vendor service
type VendorService interface {
	RegisterVendor(ctx context.Context, registerDto *dto.VendorRegisterRequestDto) (*dto.VendorResponse, *response.ErrorDetails)
	GetVendorInfo(ctx context.Context, vendorID uint) (*dto.VendorUserInfo, *response.ErrorDetails)
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

	// Hash password
	hashedPassword, err := auth.HashPassword(registerDto.Password)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to process password",
			Error:   err,
		}
	}

	// Set vendor role
	vendorRole := models.ROLE_VENDOR

	// Create user
	user := &models.User{
		Name:         registerDto.Name,
		Email:        registerDto.Email,
		Phone:        registerDto.Phone,
		PasswordHash: hashedPassword,
		RoleID:       &vendorRole,
	}

	createdUser, err := s.repository.CreateUser(user)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create user",
			Error:   err,
		}
	}

	// Generate unique vendor code
	vendorCode, err := utils.GenerateVendorCode(registerDto.CompanyName)
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate vendor code",
			Error:   err,
		}
	}

	// Create vendor profile
	vendor := &models.Vendor{
		UserID:      createdUser.ID,
		VendorCode:  vendorCode,
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
			VendorCode:  createdVendor.VendorCode,
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

// GetVendorInfo gets vendor information by ID
func (s *VendorServiceImpl) GetVendorInfo(ctx context.Context, vendorID uint) (*dto.VendorUserInfo, *response.ErrorDetails) {
	// Get vendor with user info
	vendor, err := s.repository.GetVendorByID(vendorID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &response.ErrorDetails{
				Code:    http.StatusNotFound,
				Message: "Vendor not found",
				Error:   err,
			}
		}
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get vendor information",
			Error:   err,
		}
	}

	// Convert to DTO
	vendorInfo := &dto.VendorUserInfo{
		ID:          vendor.User.ID,
		Name:        vendor.User.Name,
		Email:       vendor.User.Email,
		Phone:       vendor.User.Phone,
		VendorCode:  vendor.VendorCode,
		CompanyName: vendor.CompanyName,
		GSTNumber:   vendor.GSTNumber,
		Address:     vendor.Address,
		LogoURL:     vendor.LogoURL,
		Role:        "vendor",
	}

	return vendorInfo, nil
}
