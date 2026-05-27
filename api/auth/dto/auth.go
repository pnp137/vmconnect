package dto

import "regexp"

// LoginRequest represents the common login request for both merchant and vendor
type LoginRequest struct {
	Username string `json:"username" validate:"required"` // Can be email or mobile number
	Password string `json:"password" validate:"required"`
}

// LoginRequestDto represents the internal DTO for login processing
type LoginRequestDto struct {
	Username string
	Password string
	IsEmail  bool // true if username is email, false if mobile
}

// GenerateOTPRequest represents a request to generate an OTP for login.
type GenerateOTPRequest struct {
	Username string `json:"username" validate:"required"` // Can be email or mobile number
}

// GenerateOTPRequestDto represents the internal DTO for OTP generation.
type GenerateOTPRequestDto struct {
	Username string
	IsEmail  bool // true if username is email, false if mobile
}

// ValidateOTPRequest represents a request to validate an OTP and login.
type ValidateOTPRequest struct {
	Username string `json:"username" validate:"required"` // Can be email or mobile number
	OTP      string `json:"otp" validate:"required"`
}

// ValidateOTPRequestDto represents the internal DTO for OTP validation.
type ValidateOTPRequestDto struct {
	Username string
	OTP      string
	IsEmail  bool // true if username is email, false if mobile
}

// OTPResponse represents the OTP generation response.
type OTPResponse struct {
	Message string `json:"message"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token   string      `json:"token"`
	User    interface{} `json:"user"`
	Message string      `json:"message"`
}

// UserInfo represents basic user information
type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Role  string `json:"role"`
}

// MerchantUserInfo represents merchant-specific user info
type MerchantUserInfo struct {
	UserInfo
	MerchantID   uint   `json:"merchant_id"`
	BusinessName string `json:"business_name"`
	ShopName     string `json:"shop_name"`
	Address      string `json:"address"`
	Pincode      string `json:"pincode"`
}

// VendorUserInfo represents vendor-specific user info
type VendorUserInfo struct {
	UserInfo
	VendorID    uint   `json:"vendor_id"`
	CompanyName string `json:"company_name"`
	GSTNumber   string `json:"gst_number"`
	Address     string `json:"address"`
	LogoURL     string `json:"logo_url"`
}

// ToLoginRequestDto converts LoginRequest to LoginRequestDto
func ToLoginRequestDto(req *LoginRequest) *LoginRequestDto {
	// Determine if username is email or mobile number
	isEmail := IsEmail(req.Username)

	return &LoginRequestDto{
		Username: req.Username,
		Password: req.Password,
		IsEmail:  isEmail,
	}
}

// ToGenerateOTPRequestDto converts GenerateOTPRequest to GenerateOTPRequestDto.
func ToGenerateOTPRequestDto(req *GenerateOTPRequest) *GenerateOTPRequestDto {
	isEmail := IsEmail(req.Username)

	return &GenerateOTPRequestDto{
		Username: req.Username,
		IsEmail:  isEmail,
	}
}

// ToValidateOTPRequestDto converts ValidateOTPRequest to ValidateOTPRequestDto.
func ToValidateOTPRequestDto(req *ValidateOTPRequest) *ValidateOTPRequestDto {
	isEmail := IsEmail(req.Username)

	return &ValidateOTPRequestDto{
		Username: req.Username,
		OTP:      req.OTP,
		IsEmail:  isEmail,
	}
}

// IsEmail returns true when the username is an email address.
func IsEmail(username string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, username)
	return isEmail
}
