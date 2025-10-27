package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

type VendorMerchantRepository interface {
	GetVendorMerchants(vendorID uint) ([]models.Merchant, error)
	GetProductVisibilityForMerchant(vendorID, merchantID uint) ([]models.MerchantProductVisibility, error)
	UpdateProductVisibilityForMerchant(vendorID, merchantID uint, updates []dto.UpdateProductVisibilityRequest) error
}

type VendorMerchantRepositoryImpl struct {
	db *gorm.DB
}

func NewVendorMerchantRepository(db *system.DataSource) VendorMerchantRepository {
	return &VendorMerchantRepositoryImpl{db: db.Db}
}

// GetVendorMerchants gets all merchants linked to a vendor
func (r *VendorMerchantRepositoryImpl) GetVendorMerchants(vendorID uint) ([]models.Merchant, error) {
	var merchants []models.Merchant
	err := r.db.Table("merchants").
		Joins("JOIN vendor_merchant_mappings ON merchants.id = vendor_merchant_mappings.merchant_id").
		Preload("User").
		Where("vendor_merchant_mappings.vendor_id = ? AND vendor_merchant_mappings.is_deleted = ? AND merchants.is_deleted = ?",
			vendorID, false, false).
		Find(&merchants).Error
	return merchants, err
}

// GetProductVisibilityForMerchant gets product visibility settings for a specific merchant
func (r *VendorMerchantRepositoryImpl) GetProductVisibilityForMerchant(vendorID, merchantID uint) ([]models.MerchantProductVisibility, error) {
	var visibilities []models.MerchantProductVisibility
	err := r.db.Preload("Product.Category").
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
