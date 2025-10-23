package validator

import (
	"errors"
	"regexp"

	"linksupply.io/vmconnect/api/merchant/dto"
)

// MerchantValidator defines the interface for merchant validation
type MerchantValidator interface {
	ValidateRegisterRequest(req *dto.MerchantRegisterRequest) error
}

type MerchantValidatorImpl struct{}

// NewMerchantValidator creates a new instance of merchant validator
func NewMerchantValidator() MerchantValidator {
	return &MerchantValidatorImpl{}
}

// ValidateRegisterRequest validates the merchant registration request
func (v *MerchantValidatorImpl) ValidateRegisterRequest(req *dto.MerchantRegisterRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}

	if req.Email == "" {
		return errors.New("email is required")
	}

	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, req.Email); !matched {
		return errors.New("invalid email format")
	}

	if req.Phone == "" {
		return errors.New("phone is required")
	}
	if len(req.Phone) < 10 || len(req.Phone) > 15 {
		return errors.New("phone must be between 10 and 15 characters")
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if req.ShopName == "" {
		return errors.New("shop name is required")
	}
	if len(req.ShopName) < 2 || len(req.ShopName) > 100 {
		return errors.New("shop name must be between 2 and 100 characters")
	}

	if req.Address == "" {
		return errors.New("address is required")
	}

	if req.Pincode == "" {
		return errors.New("pincode is required")
	}
	if len(req.Pincode) != 6 {
		return errors.New("pincode must be exactly 6 characters")
	}

	// Validate pincode is numeric
	if matched, _ := regexp.MatchString(`^\d{6}$`, req.Pincode); !matched {
		return errors.New("pincode must be 6 digits")
	}

	return nil
}
