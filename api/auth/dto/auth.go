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
	ShopName string `json:"shop_name"`
	Address  string `json:"address"`
	Pincode  string `json:"pincode"`
}

// VendorUserInfo represents vendor-specific user info
type VendorUserInfo struct {
	UserInfo
	CompanyName string `json:"company_name"`
	GSTNumber   string `json:"gst_number"`
	Address     string `json:"address"`
	LogoURL     string `json:"logo_url"`
}

// ToLoginRequestDto converts LoginRequest to LoginRequestDto
func ToLoginRequestDto(req *LoginRequest) *LoginRequestDto {
	// Determine if username is email or mobile number
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	isEmail, _ := regexp.MatchString(emailRegex, req.Username)

	return &LoginRequestDto{
		Username: req.Username,
		Password: req.Password,
		IsEmail:  isEmail,
	}
}
