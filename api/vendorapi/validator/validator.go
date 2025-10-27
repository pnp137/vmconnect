package validator

import (
	"errors"
	"regexp"
	"strconv"

	"linksupply.io/vmconnect/api/vendorapi/dto"
)

// VendorValidator defines the interface for vendor validation
type VendorValidator interface {
	ValidateRegisterRequest(req *dto.VendorRegisterRequest) error
	ValidateAddProductRequest(req *dto.AddProductRequest) error
	ValidateUpdateProductRequest(req *dto.UpdateProductRequest) error
	ValidateProductQueryParam(queryParams *dto.ProductQueryParam) error
	ValidateIDParameter(idStr, paramName string) (uint, error)
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

// ValidateAddProductRequest validates the add product request
func (v *VendorValidatorImpl) ValidateAddProductRequest(req *dto.AddProductRequest) error {
	if req.Name == "" {
		return errors.New("product name is required")
	}
	if len(req.Name) < 1 || len(req.Name) > 255 {
		return errors.New("product name must be between 1 and 255 characters")
	}

	if len(req.Description) > 1000 {
		return errors.New("description cannot exceed 1000 characters")
	}

	if req.Price < 0 {
		return errors.New("price must be non-negative")
	}

	// Category validation: either category_id or category_name must be provided
	if req.CategoryID == 0 && req.CategoryName == "" {
		return errors.New("either category_id or category_name is required")
	}

	if req.CategoryName != "" && (len(req.CategoryName) < 1 || len(req.CategoryName) > 100) {
		return errors.New("category name must be between 1 and 100 characters")
	}

	if req.SKU < 1 {
		return errors.New("SKU must be greater than 0")
	}

	if len(req.ImageURL) > 500 {
		return errors.New("image URL cannot exceed 500 characters")
	}

	if req.Stock < 0 {
		return errors.New("stock must be non-negative")
	}

	return nil
}

// ValidateUpdateProductRequest validates the update product request
func (v *VendorValidatorImpl) ValidateUpdateProductRequest(req *dto.UpdateProductRequest) error {
	if req.Name != nil {
		if len(*req.Name) < 1 || len(*req.Name) > 255 {
			return errors.New("product name must be between 1 and 255 characters")
		}
	}

	if req.Description != nil {
		if len(*req.Description) > 1000 {
			return errors.New("description cannot exceed 1000 characters")
		}
	}

	if req.Price != nil {
		if *req.Price < 0 {
			return errors.New("price must be non-negative")
		}
	}

	if req.SKU != nil {
		if len(*req.SKU) > 50 {
			return errors.New("SKU cannot exceed 50 characters")
		}
	}

	if req.ImageURL != nil {
		if len(*req.ImageURL) > 500 {
			return errors.New("image URL cannot exceed 500 characters")
		}
	}

	if req.Stock != nil {
		if *req.Stock < 0 {
			return errors.New("stock must be non-negative")
		}
	}

	return nil
}

// ValidateProductQueryParam validates the product query parameters
func (v *VendorValidatorImpl) ValidateProductQueryParam(queryParams *dto.ProductQueryParam) error {
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	if queryParams.Offset < 0 {
		return errors.New("offset must be non-negative")
	}

	if queryParams.Ordering != "" {
		validOrderings := []string{"name", "price", "created_at", "updated_at"}
		isValid := false
		for _, validOrdering := range validOrderings {
			if queryParams.Ordering == validOrdering {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("ordering must be one of: name, price, created_at, updated_at")
		}
	}

	if len(queryParams.Search) > 100 {
		return errors.New("search term cannot exceed 100 characters")
	}

	return nil
}

// ValidateIDParameter validates and parses ID parameters from URL
func (v *VendorValidatorImpl) ValidateIDParameter(idStr, paramName string) (uint, error) {
	if idStr == "" {
		return 0, errors.New(paramName + " is required")
	}

	// Parse as uint
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, errors.New("invalid " + paramName + " format")
	}

	// Check if ID is valid (greater than 0)
	if id == 0 {
		return 0, errors.New(paramName + " must be greater than 0")
	}

	// Check reasonable upper limit
	if id > 999999999 {
		return 0, errors.New(paramName + " is too large")
	}

	return uint(id), nil
}
