package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// VendorRepository defines the interface for vendor repository
type VendorRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	CreateVendor(vendor *models.Vendor) (*models.Vendor, error)
	GetUserByEmailOrPhone(email, phone string) (*models.User, error)
	GetRoleByName(roleName string) (*models.Role, error)
}

type VendorRepositoryImpl struct {
	db *gorm.DB
}

// NewVendorRepository creates a new instance of vendor repository
func NewVendorRepository(db *system.DataSource) VendorRepository {
	return &VendorRepositoryImpl{db: db.Db}
}

// CreateUser creates a new user
func (r *VendorRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateVendor creates a new vendor profile
func (r *VendorRepositoryImpl) CreateVendor(vendor *models.Vendor) (*models.Vendor, error) {
	if err := r.db.Create(vendor).Error; err != nil {
		return nil, err
	}
	return vendor, nil
}

// GetUserByEmailOrPhone checks if user exists with email or phone
func (r *VendorRepositoryImpl) GetUserByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND is_deleted = ?", email, phone, false).First(&user).Error
	return &user, err
}

// GetRoleByName retrieves role by name
func (r *VendorRepositoryImpl) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ? AND is_deleted = ?", roleName, false).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
