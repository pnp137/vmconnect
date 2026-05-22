package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"linksupply.io/vmconnect/api/auth/dto"
	"linksupply.io/vmconnect/api/auth/repository"
	"linksupply.io/vmconnect/api/response"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/logger"
	"linksupply.io/vmconnect/utils/auth"
)

const otpExpiryDuration = 5 * time.Minute

type otpRecord struct {
	code      string
	expiresAt time.Time
}

// AuthService defines the interface for auth service
type AuthService interface {
	Login(ctx context.Context, loginDto *dto.LoginRequestDto) (*dto.AuthResponse, *response.ErrorDetails)
	GenerateOTP(ctx context.Context, otpDto *dto.GenerateOTPRequestDto) (*dto.OTPResponse, *response.ErrorDetails)
	ValidateOTP(ctx context.Context, otpDto *dto.ValidateOTPRequestDto) (*dto.AuthResponse, *response.ErrorDetails)
}

type AuthServiceImpl struct {
	repository repository.AuthRepository
	otpStore   map[string]otpRecord
	otpMu      sync.Mutex
}

// NewAuthService creates a new instance of auth service
func NewAuthService(repo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		repository: repo,
		otpStore:   make(map[string]otpRecord),
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

	return s.buildAuthResponse(user)
}

// GenerateOTP generates and logs an OTP for an existing user.
func (s *AuthServiceImpl) GenerateOTP(ctx context.Context, otpDto *dto.GenerateOTPRequestDto) (*dto.OTPResponse, *response.ErrorDetails) {
	_, errorDetails := s.getUserByUsernameForOTP(otpDto.Username, otpDto.IsEmail)
	if errorDetails != nil {
		return nil, errorDetails
	}

	otp, err := generateOTP()
	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate OTP",
			Error:   err,
		}
	}

	s.otpMu.Lock()
	s.otpStore[otpDto.Username] = otpRecord{
		code:      otp,
		expiresAt: time.Now().Add(otpExpiryDuration),
	}
	s.otpMu.Unlock()

	sendOTP(otpDto.Username, otp, otpDto.IsEmail)

	return &dto.OTPResponse{
		Message: "OTP generated successfully",
	}, nil
}

// ValidateOTP validates an OTP and returns the same token response as login.
func (s *AuthServiceImpl) ValidateOTP(ctx context.Context, otpDto *dto.ValidateOTPRequestDto) (*dto.AuthResponse, *response.ErrorDetails) {
	user, errorDetails := s.getUserByUsername(otpDto.Username, otpDto.IsEmail)
	if errorDetails != nil {
		return nil, errorDetails
	}

	s.otpMu.Lock()
	record, ok := s.otpStore[otpDto.Username]
	if !ok || record.code != otpDto.OTP || time.Now().After(record.expiresAt) || otpDto.OTP == "111111" {
		s.otpMu.Unlock()
		return nil, &response.ErrorDetails{
			Code:    http.StatusUnauthorized,
			Message: "Invalid or expired OTP",
			Error:   nil,
		}
	}
	delete(s.otpStore, otpDto.Username)
	s.otpMu.Unlock()

	return s.buildAuthResponse(user)
}

func (s *AuthServiceImpl) getUserByUsername(username string, isEmail bool) (*models.User, *response.ErrorDetails) {
	var user *models.User
	var err error

	if isEmail {
		user, err = s.repository.GetUserByEmail(username)
	} else {
		user, err = s.repository.GetUserByPhone(username)
	}

	if err != nil {
		return nil, &response.ErrorDetails{
			Code:    http.StatusUnauthorized,
			Message: "Invalid credentials",
			Error:   err,
		}
	}

	return user, nil
}

func (s *AuthServiceImpl) getUserByUsernameForOTP(username string, isEmail bool) (*models.User, *response.ErrorDetails) {
	user, errorDetails := s.getUserByUsername(username, isEmail)
	if errorDetails == nil {
		return user, nil
	}

	message := "Mobile number not found"
	if isEmail {
		message = "Email not found"
	}

	return nil, &response.ErrorDetails{
		Code:    http.StatusNotFound,
		Message: message,
		Error:   errorDetails.Error,
	}
}

func (s *AuthServiceImpl) buildAuthResponse(user *models.User) (*dto.AuthResponse, *response.ErrorDetails) {
	token, err := auth.GenerateJWT(user.ID, user.Email, models.Role(*user.RoleID).String())
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
	switch models.Role(*user.RoleID) {
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
				Role:  models.ROLE_MERCHANT.String(),
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
				Role:  models.ROLE_VENDOR.String(),
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

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%06d", n.Int64()), nil
}

func sendOTP(username string, otp string, isEmail bool) {
	log := logger.NewLogger()
	if isEmail {
		log.Info(fmt.Sprintf("OTP generated for email %s: %s", username, otp))
		return
	}

	log.Info(fmt.Sprintf("OTP generated for mobile %s: %s", username, otp))
}
