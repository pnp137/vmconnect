package validator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"linksupply.io/vmconnect/api/vendorapi/dto"
)

// VendorValidator defines the interface for vendor validation
type VendorValidator interface {
	ValidateRegisterRequest(req *dto.VendorRegisterRequest) error
	ValidateCreateOrderRequest(req *dto.CreateOrderRequest) error
	ValidateAddProductRequest(req *dto.AddProductRequest) error
	ValidateUpdateProductRequest(req *dto.UpdateProductRequest) error
	ValidateProductQueryParam(queryParams *dto.ProductQueryParam) error
	ValidateIDParameter(idStr, paramName string) (uint, error)
	ValidateAddMerchantRequest(req *dto.AddMerchantRequest) error
	ValidateUpdateMerchantRequest(req *dto.UpdateMerchantRequest) error
	ValidateMerchantQueryParam(queryParams *dto.MerchantQueryParam) error
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

func (v *VendorValidatorImpl) ValidateCreateOrderRequest(req *dto.CreateOrderRequest) error {
	if req == nil {
		return errors.New("request body is required")
	}

	req.OrderFor = strings.ToLower(strings.TrimSpace(req.OrderFor))
	if req.OrderFor == "" {
		return errors.New("order_for is required")
	}
	if req.OrderFor != "business" && req.OrderFor != "personal" {
		return errors.New("order_for must be either business or personal")
	}

	req.CustomerName = strings.TrimSpace(req.CustomerName)
	if req.CustomerName == "" {
		return errors.New("customer_name is required")
	}

	req.CustomerMobile = strings.TrimSpace(req.CustomerMobile)
	if req.CustomerMobile == "" {
		return errors.New("customer_mobile is required")
	}

	mobileRegex := `^[6-9]\d{9}$`
	if matched, _ := regexp.MatchString(mobileRegex, req.CustomerMobile); !matched {
		return errors.New("customer_mobile must be a valid 10 digit Indian mobile number")
	}

	req.DeliveryAddress = strings.TrimSpace(req.DeliveryAddress)
	if req.DeliveryAddress == "" && req.Location == nil {
		return errors.New("delivery_address or location is required")
	}

	if req.Location != nil {
		if req.Location.Latitude < -90 || req.Location.Latitude > 90 {
			return errors.New("location.latitude must be between -90 and 90")
		}
		if req.Location.Longitude < -180 || req.Location.Longitude > 180 {
			return errors.New("location.longitude must be between -180 and 180")
		}
	}

	if len(req.Items) == 0 {
		return errors.New("items are required")
	}

	for index, item := range req.Items {
		if item.ProductVariantID == 0 {
			return errors.New("product_variant_id is required for item " + strconv.Itoa(index+1))
		}
		if item.Quantity <= 0 {
			return errors.New("quantity must be greater than 0 for all items")
		}
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

	// Category validation: category_name is required
	if req.CategoryName == "" {
		return errors.New("category_name is required")
	}

	if len(req.CategoryName) < 1 || len(req.CategoryName) > 100 {
		return errors.New("category name must be between 1 and 100 characters")
	}

	if len(req.ImageURLs) > 10 {
		return errors.New("maximum 10 images are allowed")
	}
	for _, imageURL := range req.ImageURLs {
		if len(imageURL) > 500 {
			return errors.New("each image URL cannot exceed 500 characters")
		}
	}
	if len(req.ThumbnailURL) > 500 {
		return errors.New("thumbnail URL cannot exceed 500 characters")
	}

	if len(req.Variants) == 0 {
		return errors.New("at least one variant is required")
	}
	for _, variant := range req.Variants {
		if strings.TrimSpace(variant.Name) == "" {
			return errors.New("variant name is required")
		}
		if variant.Price < 0 {
			return errors.New("variant price must be non-negative")
		}
		if variant.Stock < 0 {
			return errors.New("variant stock must be non-negative")
		}
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

	if req.ImageURLs != nil {
		if len(*req.ImageURLs) > 10 {
			return errors.New("maximum 10 images are allowed")
		}
		for _, imageURL := range *req.ImageURLs {
			if len(imageURL) > 500 {
				return errors.New("each image URL cannot exceed 500 characters")
			}
		}
	}
	if req.ThumbnailURL != nil && len(*req.ThumbnailURL) > 500 {
		return errors.New("thumbnail URL cannot exceed 500 characters")
	}

	if req.Variants != nil {
		for _, variant := range *req.Variants {
			if strings.TrimSpace(variant.Name) == "" {
				return errors.New("variant name is required")
			}
			if variant.Price < 0 {
				return errors.New("variant price must be non-negative")
			}
			if variant.Stock < 0 {
				return errors.New("variant stock must be non-negative")
			}
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

	if queryParams.CategoryID > 0 && queryParams.CategoryName != "" {
		return errors.New("category_id and category_name cannot be used together")
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

	if len(queryParams.NameOrDescription) > 100 {
		return errors.New("search term cannot exceed 100 characters")
	}
	if len(queryParams.CategoryName) > 100 {
		return errors.New("category_name cannot exceed 100 characters")
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

// ValidateAddMerchantRequest validates a vendor-created merchant request.
func (v *VendorValidatorImpl) ValidateAddMerchantRequest(req *dto.AddMerchantRequest) error {
	if req.Name == "" {
		return errors.New("name is required")
	}
	if len(req.Name) < 2 || len(req.Name) > 100 {
		return errors.New("name must be between 2 and 100 characters")
	}

	if err := validateEmail(req.Email); err != nil {
		return err
	}

	if err := validatePhone(req.Phone); err != nil {
		return err
	}

	if req.Password == "" {
		return errors.New("password is required")
	}
	if len(req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if req.BusinessName == "" {
		return errors.New("business name is required")
	}
	if len(req.BusinessName) < 2 || len(req.BusinessName) > 100 {
		return errors.New("business name must be between 2 and 100 characters")
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

	if err := validatePincode(req.Pincode); err != nil {
		return err
	}

	return nil
}

// ValidateUpdateMerchantRequest validates merchant updates made by a vendor.
func (v *VendorValidatorImpl) ValidateUpdateMerchantRequest(req *dto.UpdateMerchantRequest) error {
	if req.Name != nil && (len(*req.Name) < 2 || len(*req.Name) > 100) {
		return errors.New("name must be between 2 and 100 characters")
	}

	if req.Email != nil {
		if err := validateEmail(*req.Email); err != nil {
			return err
		}
	}

	if req.Phone != nil {
		if err := validatePhone(*req.Phone); err != nil {
			return err
		}
	}

	if req.Password != nil && len(*req.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	if req.BusinessName != nil && (len(*req.BusinessName) < 2 || len(*req.BusinessName) > 100) {
		return errors.New("business name must be between 2 and 100 characters")
	}

	if req.ShopName != nil && (len(*req.ShopName) < 2 || len(*req.ShopName) > 100) {
		return errors.New("shop name must be between 2 and 100 characters")
	}

	if req.Address != nil && *req.Address == "" {
		return errors.New("address cannot be empty")
	}

	if req.Pincode != nil {
		if err := validatePincode(*req.Pincode); err != nil {
			return err
		}
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		return errors.New("invalid email format")
	}

	return nil
}

func validatePhone(phone string) error {
	if phone == "" {
		return errors.New("phone is required")
	}
	if len(phone) < 10 || len(phone) > 15 {
		return errors.New("phone must be between 10 and 15 characters")
	}

	return nil
}

// ValidateMerchantQueryParam validates the merchant query parameters
func (v *VendorValidatorImpl) ValidateMerchantQueryParam(queryParams *dto.MerchantQueryParam) error {
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	if queryParams.Offset < 0 {
		return errors.New("offset must be non-negative")
	}

	if queryParams.Ordering != "" {
		validOrderings := []string{"owner_name", "business_name", "shop_name", "created_at"}
		isValid := false
		for _, validOrdering := range validOrderings {
			if queryParams.Ordering == validOrdering {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("ordering must be one of: owner_name, business_name, shop_name, created_at")
		}
	}

	if len(queryParams.BusinessName) > 100 {
		return errors.New("business_name cannot exceed 100 characters")
	}
	if len(queryParams.ShopName) > 100 {
		return errors.New("shop_name cannot exceed 100 characters")
	}
	if len(queryParams.Pincode) > 10 {
		return errors.New("pincode cannot exceed 10 characters")
	}
	if len(queryParams.OwnerName) > 100 {
		return errors.New("owner_name cannot exceed 100 characters")
	}

	return nil
}

func validatePincode(pincode string) error {
	if pincode == "" {
		return errors.New("pincode is required")
	}
	if len(pincode) != 6 {
		return errors.New("pincode must be exactly 6 characters")
	}

	return nil
}
