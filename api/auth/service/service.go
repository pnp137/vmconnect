package service

import (
	"context"
	"net/http"

	"linksupply.io/vmconnect/api/auth/dto"
	"linksupply.io/vmconnect/api/auth/repository"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/utils/auth"
)

// AuthService defines the interface for auth service
type AuthService interface {
	Login(ctx context.Context, loginDto *dto.LoginRequestDto) (*dto.AuthResponse, *response.ErrorDetails)
}

type AuthServiceImpl struct {
	repository repository.AuthRepository
}

// NewAuthService creates a new instance of auth service
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		repository: repo,
	}
}

// Login handles the login logic for both merchants and vendors
func (s *AuthServiceImpl) Login(ctx context.Context, loginDto *dto.LoginRequestDto) (*dto.AuthResponse, *response.ErrorDetails) {
	var user *models.User
	var err error

	// Find user by email or phone based on IsEmail flag
	if loginDto.IsEmail {
		user, err = s.repository.GetUserByEmail(loginDto.Username)
	} else {
		user, err = s.repository.GetUserByPhone(loginDto.Username)
	}

	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
			Error:   err,
		}
	}

	// Validate password
	if err := auth.CheckPassword(loginDto.Password, user.PasswordHash); err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
			Error:   err,
		}
	}

	// Generate JWT token
	token, err := auth.GenerateJWT(user.ID, user.Email, user.RoleID.String())
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
			Error:   err,
		}
	}

	var userInfo interface{}
	var message string

	// Get user-specific profile based on role
	switch *user.RoleID {
	case models.ROLE_MERCHANT:
		merchant, err := s.repository.GetMerchantByUserID(user.ID)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Merchant profile not found",
				Error:   err,
			}
		}
		userInfo = &dto.MerchantUserInfo{
			UserInfo: dto.UserInfo{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Phone: user.Phone,
				Role:  user.RoleID.String(),
			},
			ShopName: merchant.ShopName,
			Address:  merchant.Address,
			Pincode:  merchant.Pincode,
		}
		message = "Merchant login successful"

	case models.ROLE_VENDOR:
		vendor, err := s.repository.GetVendorByUserID(user.ID)
		if err != nil {
			return nil, &response.ErrorDetails{
				Code:    http.StatusInternalServerError,
				Message: "Vendor profile not found",
				Error:   err,
			}
		}
		userInfo = &dto.VendorUserInfo{
			UserInfo: dto.UserInfo{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Phone: user.Phone,
				Role:  user.RoleID.String(),
			},
			CompanyName: vendor.CompanyName,
			GSTNumber:   vendor.GSTNumber,
			Address:     vendor.Address,
			LogoURL:     vendor.LogoURL,
		}
		message = "Vendor login successful"

	default:
		return nil, &response.ErrorDetails{
			Code:    http.StatusForbidden,
			Message: "Access denied: Invalid user role",
			Error:   nil,
		}
	}

	return &dto.AuthResponse{
		Token:   token,
		User:    userInfo,
		Message: message,
	}, nil
}
