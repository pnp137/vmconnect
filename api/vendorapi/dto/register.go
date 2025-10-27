package dto

// VendorRegisterRequest represents the vendor registration request
type VendorRegisterRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Email       string `json:"email" validate:"required,email"`
	Phone       string `json:"phone" validate:"required,min=10,max=15"`
	Password    string `json:"password" validate:"required,min=6"`
	CompanyName string `json:"company_name" validate:"required,min=2,max=100"`
	GSTNumber   string `json:"gst_number" validate:"required,min=15,max=15"`
	Address     string `json:"address" validate:"required"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// VendorRegisterRequestDto represents the internal DTO for vendor registration
type VendorRegisterRequestDto struct {
	Name        string
	Email       string
	Phone       string
	Password    string
	CompanyName string
	GSTNumber   string
	Address     string
	LogoURL     string
}

// VendorResponse represents the vendor registration response
type VendorResponse struct {
	User    VendorUserInfo `json:"user"`
	Message string         `json:"message"`
}

// VendorUserInfo represents the vendor user info in response
type VendorUserInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	VendorCode  string `json:"vendor_code"`
	CompanyName string `json:"company_name"`
	GSTNumber   string `json:"gst_number"`
	Address     string `json:"address"`
	LogoURL     string `json:"logo_url"`
	Role        string `json:"role"`
}

// ToVendorRegisterRequestDto converts VendorRegisterRequest to VendorRegisterRequestDto
func ToVendorRegisterRequestDto(req *VendorRegisterRequest) *VendorRegisterRequestDto {
	return &VendorRegisterRequestDto{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Password:    req.Password,
		CompanyName: req.CompanyName,
		GSTNumber:   req.GSTNumber,
		Address:     req.Address,
		LogoURL:     req.LogoURL,
	}
}
