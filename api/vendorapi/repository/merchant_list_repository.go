package repository

import (
	"strings"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

// VendorMerchantListRepository handles paginated listing and search of vendor merchants.
type VendorMerchantListRepository interface {
	GetVendorMerchants(vendorID uint, queryParams *dto.MerchantQueryParam) ([]models.Merchant, int, error)
}

type VendorMerchantListRepositoryImpl struct {
	db *gorm.DB
}

func NewVendorMerchantListRepository(db *system.DataSource) VendorMerchantListRepository {
	return &VendorMerchantListRepositoryImpl{db: db.Db}
}

// GetVendorMerchants gets merchants linked to a vendor with pagination and optional filters.
func (r *VendorMerchantListRepositoryImpl) GetVendorMerchants(vendorID uint, queryParams *dto.MerchantQueryParam) ([]models.Merchant, int, error) {
	var merchants []models.Merchant
	var totalCount int64

	countQuery := r.buildFilteredQuery(vendorID, queryParams)
	if err := countQuery.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	findQuery := r.buildFilteredQuery(vendorID, queryParams).Preload("Owner")
	findQuery = r.applyOrdering(findQuery, queryParams.Ordering)

	if err := findQuery.Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&merchants).Error; err != nil {
		return nil, 0, err
	}

	return merchants, int(totalCount), nil
}

func (r *VendorMerchantListRepositoryImpl) buildFilteredQuery(vendorID uint, params *dto.MerchantQueryParam) *gorm.DB {
	query := r.db.Model(&models.Merchant{}).
		Joins("JOIN vendor_merchant_mappings ON merchants.id = vendor_merchant_mappings.merchant_id").
		Joins("JOIN users ON users.id = merchants.user_id").
		Where("vendor_merchant_mappings.vendor_id = ? AND vendor_merchant_mappings.is_deleted = ? AND merchants.is_deleted = ?",
			vendorID, false, false)

	return r.applyMerchantFilters(query, params)
}

func (r *VendorMerchantListRepositoryImpl) applyOrdering(query *gorm.DB, ordering string) *gorm.DB {
	switch ordering {
	case "owner_name":
		return query.Order("users.name ASC")
	case "business_name":
		return query.Order("merchants.business_name ASC")
	case "shop_name":
		return query.Order("merchants.shop_name ASC")
	case "created_at":
		return query.Order("merchants.created_at DESC")
	default:
		return query.Order("merchants.created_at DESC")
	}
}

func (r *VendorMerchantListRepositoryImpl) applyMerchantFilters(query *gorm.DB, params *dto.MerchantQueryParam) *gorm.DB {
	if businessName := strings.TrimSpace(params.BusinessName); businessName != "" {
		query = query.Where("merchants.business_name LIKE ?", "%"+businessName+"%")
	}
	if shopName := strings.TrimSpace(params.ShopName); shopName != "" {
		query = query.Where("merchants.shop_name LIKE ?", "%"+shopName+"%")
	}
	if pincode := strings.TrimSpace(params.Pincode); pincode != "" {
		query = query.Where("merchants.pincode LIKE ?", "%"+pincode+"%")
	}
	if ownerName := strings.TrimSpace(params.OwnerName); ownerName != "" {
		query = query.Where("users.name LIKE ?", "%"+ownerName+"%")
	}
	if params.MerchantID > 0 {
		query = query.Where("merchants.id = ?", params.MerchantID)
	}
	return query
}
