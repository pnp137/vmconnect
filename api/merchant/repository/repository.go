package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// MerchantRepository defines the interface for merchant repository
type MerchantRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	CreateMerchant(merchant *models.Merchant) (*models.Merchant, error)
	GetUserByEmailOrPhone(email, phone string) (*models.User, error)
	GetRoleByName(roleName string) (*models.Role, error)
}

type MerchantRepositoryImpl struct {
	db *gorm.DB
}

// NewMerchantRepository creates a new instance of merchant repository
func NewMerchantRepository(db *system.DataSource) MerchantRepository {
	return &MerchantRepositoryImpl{db: db.Db}
}

// CreateUser creates a new user
func (r *MerchantRepositoryImpl) CreateUser(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateMerchant creates a new merchant profile
func (r *MerchantRepositoryImpl) CreateMerchant(merchant *models.Merchant) (*models.Merchant, error) {
	if err := r.db.Create(merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}

// GetUserByEmailOrPhone checks if user exists with email or phone
func (r *MerchantRepositoryImpl) GetUserByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND is_deleted = ?", email, phone, false).First(&user).Error
	return &user, err
}

// GetRoleByName retrieves role by name
func (r *MerchantRepositoryImpl) GetRoleByName(roleName string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ? AND is_deleted = ?", roleName, false).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
