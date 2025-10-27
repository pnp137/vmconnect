package dto

// MerchantRegisterRequest represents the merchant registration request
type MerchantRegisterRequest struct {
	Name       string `json:"name" validate:"required,min=2,max=100"`
	Email      string `json:"email" validate:"required,email"`
	Phone      string `json:"phone" validate:"required,min=10,max=15"`
	Password   string `json:"password" validate:"required,min=6"`
	ShopName   string `json:"shop_name" validate:"required,min=2,max=100"`
	Address    string `json:"address" validate:"required"`
	Pincode    string `json:"pincode" validate:"required,min=6,max=6"`
	VendorCode string `json:"vendor_code,omitempty" validate:"omitempty,min=7,max=7"`
}

// MerchantRegisterRequestDto represents the internal DTO for merchant registration
type MerchantRegisterRequestDto struct {
	Name       string
	Email      string
	Phone      string
	Password   string
	ShopName   string
	Address    string
	Pincode    string
	VendorCode string
}

// MerchantResponse represents the merchant registration response
type MerchantResponse struct {
	User    MerchantUserInfo `json:"user"`
	Message string           `json:"message"`
}

// MerchantUserInfo represents the merchant user info in response
type MerchantUserInfo struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	ShopName string `json:"shop_name"`
	Address  string `json:"address"`
	Pincode  string `json:"pincode"`
	Role     string `json:"role"`
}

// ToMerchantRegisterRequestDto converts MerchantRegisterRequest to MerchantRegisterRequestDto
func ToMerchantRegisterRequestDto(req *MerchantRegisterRequest) *MerchantRegisterRequestDto {
	return &MerchantRegisterRequestDto{
		Name:       req.Name,
		Email:      req.Email,
		Phone:      req.Phone,
		Password:   req.Password,
		ShopName:   req.ShopName,
		Address:    req.Address,
		Pincode:    req.Pincode,
		VendorCode: req.VendorCode,
	}
}
