package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// AuthRepository defines the interface for auth repository
type AuthRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByPhone(phone string) (*models.User, error)
	GetMerchantByUserID(userID uint) (*models.Merchant, error)
	GetVendorByUserID(userID uint) (*models.Vendor, error)
}

type AuthRepositoryImpl struct {
	db *gorm.DB
}

// NewAuthRepository creates a new instance of auth repository
func NewAuthRepository(db *system.DataSource) AuthRepository {
	return &AuthRepositoryImpl{db: db.Db}
}

// GetUserByEmail retrieves user by email
func (r *AuthRepositoryImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? AND is_deleted = ?", email, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByPhone retrieves user by phone number
func (r *AuthRepositoryImpl) GetUserByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ? AND is_deleted = ?", phone, false).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetMerchantByUserID retrieves merchant profile by user ID
func (r *AuthRepositoryImpl) GetMerchantByUserID(userID uint) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// // GetMerchantByUserID retrieves merchant profile by user ID
// func (r *AuthRepositoryImpl) GetMerchantByUserID(userID uint) (*models.Merchant, error) {
// 	var merchant models.Merchant
// 	err := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).First(&merchant).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &merchant, nil
// }

// GetVendorByUserID retrieves vendor profile by user ID
func (r *AuthRepositoryImpl) GetVendorByUserID(userID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("user_id = ? AND is_deleted = ?", userID, false).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}
