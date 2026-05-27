package validator

import (
	"errors"
	"regexp"
	"strconv"

	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/database/models"
)

// MerchantValidator defines the interface for merchant validation
type MerchantValidator interface {
	ValidateRegisterRequest(req *dto.MerchantRegisterRequest) error
	ValidateOrderQueryParam(queryParams *dto.OrderQueryParam) error
	ValidateAddVendorRequest(req *dto.AddVendorRequest) error
	ValidateCartItemRequest(req *dto.CartItemRequest) error
	ValidateUpdateOrderStatusRequest(req *dto.UpdateOrderStatusRequest) error
	ValidateProductQueryParam(queryParams *dto.ProductQueryParam) error
	ValidateIDParameter(idStr, paramName string) (uint, error)
	ValidateOrderStatusForCartOperation(orderStatus *models.OrderStatus) error
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

// ValidateOrderQueryParam validates the order query parameters
func (v *MerchantValidatorImpl) ValidateOrderQueryParam(queryParams *dto.OrderQueryParam) error {
	// Validate limit
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	// Validate offset
	if queryParams.Offset < 0 {
		return errors.New("offset must be non-negative")
	}

	// Validate ordering
	if queryParams.Ordering != "" {
		validOrderings := []string{"created_at", "updated_at", "total_amount"}
		isValid := false
		for _, validOrdering := range validOrderings {
			if queryParams.Ordering == validOrdering {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("ordering must be one of: created_at, updated_at, total_amount")
		}
	}

	// Validate order status values
	if len(queryParams.OrderStatus) > 0 {
		validStatuses := []string{
			"IN_CART", "PLACED", "CONFIRMED", "INVOICED", "PARTIALLY_PAID",
			"PAID", "READY_TO_SHIP", "SHIPPED", "DELIVERED", "COMPLETED",
			"CANCELLED", "RETURNED",
		}

		for _, status := range queryParams.OrderStatus {
			status = regexp.MustCompile(`\s+`).ReplaceAllString(status, "") // Remove whitespace
			isValid := false
			for _, validStatus := range validStatuses {
				if status == validStatus {
					isValid = true
					break
				}
			}
			if !isValid {
				return errors.New("invalid order status: " + status + ". Valid statuses are: IN_CART, PLACED, CONFIRMED, INVOICED, PARTIALLY_PAID, PAID, READY_TO_SHIP, SHIPPED, DELIVERED, COMPLETED, CANCELLED, RETURNED")
			}
		}
	}

	return nil
}

// ValidateAddVendorRequest validates the add vendor request
func (v *MerchantValidatorImpl) ValidateAddVendorRequest(req *dto.AddVendorRequest) error {
	if req.VendorCode == "" {
		return errors.New("vendor code is required")
	}

	// Validate vendor code format (7 characters: 3 letters + 4 digits)
	if len(req.VendorCode) != 7 {
		return errors.New("vendor code must be exactly 7 characters")
	}

	// Check if vendor code matches pattern (3 uppercase letters + 4 digits)
	vendorCodeRegex := `^[A-Z]{3}[0-9]{4}$`
	if matched, _ := regexp.MatchString(vendorCodeRegex, req.VendorCode); !matched {
		return errors.New("vendor code must be in format ABC1234 (3 uppercase letters followed by 4 digits)")
	}

	return nil
}

// ValidateCartItemRequest validates the cart item request
func (v *MerchantValidatorImpl) ValidateCartItemRequest(req *dto.CartItemRequest) error {
	if req.ProductID == 0 {
		return errors.New("product ID is required")
	}

	if req.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	if req.Quantity > 1000 {
		return errors.New("quantity cannot exceed 1000")
	}

	return nil
}

// ValidateUpdateOrderStatusRequest validates the update order status request
func (v *MerchantValidatorImpl) ValidateUpdateOrderStatusRequest(req *dto.UpdateOrderStatusRequest) error {
	if req.Status == "" {
		return errors.New("status is required")
	}

	// Validate order status values
	validStatuses := []string{
		"IN_CART", "PLACED", "CONFIRMED", "INVOICED", "PARTIALLY_PAID",
		"PAID", "READY_TO_SHIP", "SHIPPED", "DELIVERED", "COMPLETED",
		"CANCELLED", "RETURNED",
	}

	isValid := false
	for _, validStatus := range validStatuses {
		if req.Status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return errors.New("invalid order status: " + req.Status + ". Valid statuses are: IN_CART, PLACED, CONFIRMED, INVOICED, PARTIALLY_PAID, PAID, READY_TO_SHIP, SHIPPED, DELIVERED, COMPLETED, CANCELLED, RETURNED")
	}

	// Validate remarks length if provided
	if req.Remarks != "" && len(req.Remarks) > 500 {
		return errors.New("remarks cannot exceed 500 characters")
	}

	return nil
}

// ValidateProductQueryParam validates the product query parameters
func (v *MerchantValidatorImpl) ValidateProductQueryParam(queryParams *dto.ProductQueryParam) error {
	// Validate category ID
	if queryParams.CategoryID == 0 {
		return errors.New("category ID is required")
	}

	// Validate limit
	if queryParams.Limit < 1 || queryParams.Limit > 100 {
		return errors.New("limit must be between 1 and 100")
	}

	// Validate offset
	if queryParams.Offset < 0 {
		return errors.New("offset must be non-negative")
	}

	// Validate ordering
	if queryParams.Ordering != "" {
		validOrderings := []string{"name", "price", "created_at"}
		isValid := false
		for _, validOrdering := range validOrderings {
			if queryParams.Ordering == validOrdering {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("ordering must be one of: name, price, created_at")
		}
	}

	// Validate state
	if queryParams.State != "" {
		validStates := []string{"active", "featured"}
		isValid := false
		for _, validState := range validStates {
			if queryParams.State == validState {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("state must be one of: active, featured")
		}
	}

	return nil
}

// ValidateIDParameter validates and parses ID parameters from URL
func (v *MerchantValidatorImpl) ValidateIDParameter(idStr, paramName string) (uint, error) {
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

// ValidateOrderStatusForCartOperation validates that order is in CART status for cart operations
func (v *MerchantValidatorImpl) ValidateOrderStatusForCartOperation(orderStatus *models.OrderStatus) error {
	if orderStatus == nil {
		return errors.New("order status is required")
	}

	if *orderStatus != models.ORDER_PLACED {
		return errors.New("cart operations are only allowed when order is in PLACED status")
	}

	return nil
}
