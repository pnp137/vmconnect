package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

type VendorMerchantRepository interface {
	GetUserByEmailOrPhone(email, phone string) (*models.User, error)
	GetUserByEmailOrPhoneExcludingUser(email, phone string, userID uint) (*models.User, error)
	GetVendorByID(vendorID uint) (*models.Vendor, error)
	GetVendorMerchant(vendorID, merchantID uint) (*models.Merchant, error)
	CreateMerchantForVendor(user *models.User, merchant *models.Merchant, mapping *models.VendorMerchantMapping) (*models.Merchant, error)
	UpdateMerchantForVendor(merchant *models.Merchant, userUpdates map[string]interface{}, merchantUpdates map[string]interface{}) (*models.Merchant, error)
	DeleteMerchantForVendor(vendorID, merchantID uint) error
	GetProductVisibilityForMerchant(vendorID, merchantID uint) ([]models.MerchantProductVisibility, error)
	UpdateProductVisibilityForMerchant(vendorID, merchantID uint, updates []dto.UpdateProductVisibilityRequest) error
}

type VendorMerchantRepositoryImpl struct {
	db *gorm.DB
}

func NewVendorMerchantRepository(db *system.DataSource) VendorMerchantRepository {
	return &VendorMerchantRepositoryImpl{db: db.Db}
}

// GetUserByEmailOrPhone checks if a non-deleted user exists with email or phone.
func (r *VendorMerchantRepositoryImpl) GetUserByEmailOrPhone(email, phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND is_deleted = ?", email, phone, false).First(&user).Error
	return &user, err
}

// GetUserByEmailOrPhoneExcludingUser checks uniqueness while excluding the current user.
func (r *VendorMerchantRepositoryImpl) GetUserByEmailOrPhoneExcludingUser(email, phone string, userID uint) (*models.User, error) {
	var user models.User
	err := r.db.Where("(email = ? OR phone = ?) AND id <> ? AND is_deleted = ?", email, phone, userID, false).First(&user).Error
	return &user, err
}

// GetVendorByID gets a vendor by ID.
func (r *VendorMerchantRepositoryImpl) GetVendorByID(vendorID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Where("id = ? AND is_deleted = ?", vendorID, false).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// GetVendorMerchant gets a merchant that belongs to the given vendor.
func (r *VendorMerchantRepositoryImpl) GetVendorMerchant(vendorID, merchantID uint) (*models.Merchant, error) {
	var merchant models.Merchant
	err := r.db.Table("merchants").
		Joins("JOIN vendor_merchant_mappings ON merchants.id = vendor_merchant_mappings.merchant_id").
		Preload("Owner").
		Where("vendor_merchant_mappings.vendor_id = ? AND merchants.id = ? AND vendor_merchant_mappings.is_deleted = ? AND merchants.is_deleted = ?",
			vendorID, merchantID, false, false).
		First(&merchant).Error
	if err != nil {
		return nil, err
	}
	return &merchant, nil
}

// CreateMerchantForVendor creates user, merchant, and vendor mapping atomically.
func (r *VendorMerchantRepositoryImpl) CreateMerchantForVendor(user *models.User, merchant *models.Merchant, mapping *models.VendorMerchantMapping) (*models.Merchant, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		merchant.OwnerID = user.ID
		if err := tx.Create(merchant).Error; err != nil {
			return err
		}

		mapping.MerchantID = merchant.ID
		if err := tx.Create(mapping).Error; err != nil {
			return err
		}

		return tx.Preload("Owner").First(merchant, merchant.ID).Error
	})
	if err != nil {
		return nil, err
	}

	return merchant, nil
}

// UpdateMerchantForVendor updates user and merchant records atomically.
func (r *VendorMerchantRepositoryImpl) UpdateMerchantForVendor(merchant *models.Merchant, userUpdates map[string]interface{}, merchantUpdates map[string]interface{}) (*models.Merchant, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if len(userUpdates) > 0 {
			if err := tx.Model(&models.User{}).Where("id = ? AND is_deleted = ?", merchant.OwnerID, false).Updates(userUpdates).Error; err != nil {
				return err
			}
		}

		if len(merchantUpdates) > 0 {
			if err := tx.Model(&models.Merchant{}).Where("id = ? AND is_deleted = ?", merchant.ID, false).Updates(merchantUpdates).Error; err != nil {
				return err
			}
		}

		return tx.Preload("Owner").Where("id = ? AND is_deleted = ?", merchant.ID, false).First(merchant).Error
	})
	if err != nil {
		return nil, err
	}

	return merchant, nil
}

// DeleteMerchantForVendor soft deletes the vendor-created merchant relationship and merchant.
func (r *VendorMerchantRepositoryImpl) DeleteMerchantForVendor(vendorID, merchantID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var merchant models.Merchant
		if err := tx.Table("merchants").
			Joins("JOIN vendor_merchant_mappings ON merchants.id = vendor_merchant_mappings.merchant_id").
			Where("vendor_merchant_mappings.vendor_id = ? AND merchants.id = ? AND vendor_merchant_mappings.is_deleted = ? AND merchants.is_deleted = ?",
				vendorID, merchantID, false, false).
			First(&merchant).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.VendorMerchantMapping{}).
			Where("vendor_id = ? AND merchant_id = ? AND is_deleted = ?", vendorID, merchantID, false).
			Updates(map[string]interface{}{"is_deleted": true, "status": "inactive"}).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Merchant{}).Where("id = ?", merchant.ID).Update("is_deleted", true).Error; err != nil {
			return err
		}

		return tx.Model(&models.User{}).Where("id = ?", merchant.OwnerID).Update("is_deleted", true).Error
	})
}

// GetProductVisibilityForMerchant gets product visibility settings for a specific merchant
func (r *VendorMerchantRepositoryImpl) GetProductVisibilityForMerchant(vendorID, merchantID uint) ([]models.MerchantProductVisibility, error) {
	var visibilities []models.MerchantProductVisibility
	err := r.db.Preload("Product.Category").
		Preload("Product.Variants", "is_active = ?", true).
		Where("vendor_id = ? AND merchant_id = ? AND is_deleted = ?", vendorID, merchantID, false).
		Find(&visibilities).Error
	return visibilities, err
}

// UpdateProductVisibilityForMerchant updates product visibility settings for a merchant
func (r *VendorMerchantRepositoryImpl) UpdateProductVisibilityForMerchant(vendorID, merchantID uint, updates []dto.UpdateProductVisibilityRequest) error {
	for _, update := range updates {
		// Check if visibility record exists
		var visibility models.MerchantProductVisibility
		err := r.db.Where("vendor_id = ? AND merchant_id = ? AND product_id = ? AND is_deleted = ?",
			vendorID, merchantID, update.ProductID, false).First(&visibility).Error

		if err == nil {
			// Update existing record - toggle visibility via is_deleted
			err = r.db.Model(&visibility).Update("is_deleted", !update.Visible).Error
			if err != nil {
				return err
			}
		} else if err == gorm.ErrRecordNotFound {
			// Create new record
			visibility = models.MerchantProductVisibility{
				VendorID:   vendorID,
				MerchantID: merchantID,
				ProductID:  update.ProductID,
				IsDeleted:  !update.Visible, // If visible=true, then is_deleted=false
			}
			err = r.db.Create(&visibility).Error
			if err != nil {
				return err
			}
		} else {
			// Some other error
			return err
		}
	}
	return nil
}
