package validator

import (
	"errors"
	"regexp"

	"linksupply.io/vmconnect/api/auth/dto"
)

// AuthValidator defines the interface for auth validation
type AuthValidator interface {
	ValidateLoginRequest(req *dto.LoginRequest) error
	ValidateGenerateOTPRequest(req *dto.GenerateOTPRequest) error
	ValidateValidateOTPRequest(req *dto.ValidateOTPRequest) error
}

type AuthValidatorImpl struct{}

// NewAuthValidator creates a new instance of auth validator
func NewAuthValidator() AuthValidator {
	return &AuthValidatorImpl{}
}

// ValidateLoginRequest validates the login request
func (v *AuthValidatorImpl) ValidateLoginRequest(req *dto.LoginRequest) error {
	if req.Password == "" {
		return errors.New("password is required")
	}

	if err := validateUsername(req.Username); err != nil {
		return err
	}

	// Validate password length
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

// ValidateGenerateOTPRequest validates the OTP generation request.
func (v *AuthValidatorImpl) ValidateGenerateOTPRequest(req *dto.GenerateOTPRequest) error {
	return validateUsername(req.Username)
}

// ValidateValidateOTPRequest validates the OTP validation request.
func (v *AuthValidatorImpl) ValidateValidateOTPRequest(req *dto.ValidateOTPRequest) error {
	if err := validateUsername(req.Username); err != nil {
		return err
	}

	if req.OTP == "" {
		return errors.New("otp is required")
	}

	otpRegex := `^[0-9]{6}$`
	isOTP, _ := regexp.MatchString(otpRegex, req.OTP)
	if !isOTP {
		return errors.New("otp must be a 6 digit number")
	}

	return nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	mobileRegex := `^[0-9]{10,15}$`

	isEmail, _ := regexp.MatchString(emailRegex, username)
	isMobile, _ := regexp.MatchString(mobileRegex, username)

	if !isEmail && !isMobile {
		return errors.New("username must be a valid email address or mobile number")
	}

	return nil
}
