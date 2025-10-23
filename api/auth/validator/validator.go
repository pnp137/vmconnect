package validator

import (
	"errors"
	"regexp"

	"linksupply.io/vmconnect/api/auth/dto"
)

// AuthValidator defines the interface for auth validation
type AuthValidator interface {
	ValidateLoginRequest(req *dto.LoginRequest) error
}

type AuthValidatorImpl struct{}

// NewAuthValidator creates a new instance of auth validator
func NewAuthValidator() AuthValidator {
	return &AuthValidatorImpl{}
}

// ValidateLoginRequest validates the login request
func (v *AuthValidatorImpl) ValidateLoginRequest(req *dto.LoginRequest) error {
	if req.Username == "" {
		return errors.New("username is required")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}

	// Validate username format (email or mobile number)
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	mobileRegex := `^[0-9]{10,15}$`

	isEmail, _ := regexp.MatchString(emailRegex, req.Username)
	isMobile, _ := regexp.MatchString(mobileRegex, req.Username)

	if !isEmail && !isMobile {
		return errors.New("username must be a valid email address or mobile number")
	}

	// Validate password length
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}
