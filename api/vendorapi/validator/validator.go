package validator

import (
	"errors"
	"regexp"

	"linksupply.io/vmconnect/api/vendorapi/dto"
)

// VendorValidator defines the interface for vendor validation
type VendorValidator interface {
	ValidateRegisterRequest(req *dto.VendorRegisterRequest) error
}

type VendorValidatorImpl struct{}

// NewVendorValidator creates a new instance of vendor validator
func NewVendorValidator() VendorValidator {
	return &VendorValidatorImpl{}
}

// ValidateRegisterRequest validates the vendor registration request
func (v *VendorValidatorImpl) ValidateRegisterRequest(req *dto.VendorRegisterRequest) error {
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

	if req.CompanyName == "" {
		return errors.New("company name is required")
	}
	if len(req.CompanyName) < 2 || len(req.CompanyName) > 100 {
		return errors.New("company name must be between 2 and 100 characters")
	}

	if req.GSTNumber == "" {
		return errors.New("GST number is required")
	}
	if len(req.GSTNumber) != 15 {
		return errors.New("GST number must be exactly 15 characters")
	}

	// Validate GST number format (basic validation)
	gstRegex := `^[0-9]{2}[A-Z]{5}[0-9]{4}[A-Z]{1}[1-9A-Z]{1}Z[0-9A-Z]{1}$`
	if matched, _ := regexp.MatchString(gstRegex, req.GSTNumber); !matched {
		return errors.New("invalid GST number format")
	}

	if req.Address == "" {
		return errors.New("address is required")
	}

	return nil
}
