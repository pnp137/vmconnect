package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

type VendorProductRepository interface {
	CreateProduct(vendorID uint, product *models.Product) (*models.Product, error)
	GetProductByID(vendorID, productID uint) (*models.Product, error)
	UpdateProduct(vendorID, productID uint, updates map[string]interface{}) (*models.Product, error)
	DeleteProduct(vendorID, productID uint) error
	GetVendorProducts(vendorID uint, queryParams *dto.ProductQueryParam) ([]models.Product, int, error)
	GetVendorCategories(vendorID uint) ([]models.Category, error)
	GetOrCreateCategory(vendorID uint, categoryName string) (*models.Category, error)
}

type VendorProductRepositoryImpl struct {
	db *gorm.DB
}

func NewVendorProductRepository(db *system.DataSource) VendorProductRepository {
	return &VendorProductRepositoryImpl{db: db.Db}
}

// CreateProduct creates a new product for a vendor
func (r *VendorProductRepositoryImpl) CreateProduct(vendorID uint, product *models.Product) (*models.Product, error) {
	product.VendorID = vendorID
	err := r.db.Create(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

// GetProductByID gets a product by ID for a specific vendor
func (r *VendorProductRepositoryImpl) GetProductByID(vendorID, productID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").
		Where("id = ? AND vendor_id = ? AND is_deleted = ?", productID, vendorID, false).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct updates a product for a specific vendor
func (r *VendorProductRepositoryImpl) UpdateProduct(vendorID, productID uint, updates map[string]interface{}) (*models.Product, error) {
	// First verify the product belongs to the vendor
	var product models.Product
	err := r.db.Where("id = ? AND vendor_id = ? AND is_deleted = ?", productID, vendorID, false).First(&product).Error
	if err != nil {
		return nil, err
	}

	// Update the product
	err = r.db.Model(&models.Product{}).
		Where("id = ? AND vendor_id = ? AND is_deleted = ?", productID, vendorID, false).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Get the updated product with preloaded category
	err = r.db.Preload("Category").
		Where("id = ? AND vendor_id = ? AND is_deleted = ?", productID, vendorID, false).
		First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// DeleteProduct hard deletes a product for a specific vendor
func (r *VendorProductRepositoryImpl) DeleteProduct(vendorID, productID uint) error {
	// First verify the product belongs to the vendor
	var product models.Product
	err := r.db.Where("id = ? AND vendor_id = ? AND is_deleted = ?", productID, vendorID, false).First(&product).Error
	if err != nil {
		return err
	}

	// Hard delete the product
	err = r.db.Delete(&product).Error
	if err != nil {
		return err
	}

	return nil
}

// GetVendorProducts gets products for a vendor with filtering and pagination
func (r *VendorProductRepositoryImpl) GetVendorProducts(vendorID uint, queryParams *dto.ProductQueryParam) ([]models.Product, int, error) {
	var products []models.Product
	var totalCount int64

	// Build base query
	query := r.db.Model(&models.Product{}).Preload("Category").
		Where("vendor_id = ? AND is_deleted = ?", vendorID, false)

	// Apply filters
	if queryParams.CategoryID > 0 {
		query = query.Where("category_id = ?", queryParams.CategoryID)
	}

	if queryParams.IsActive != nil {
		query = query.Where("is_active = ?", *queryParams.IsActive)
	}

	if queryParams.IsFeatured != nil {
		query = query.Where("is_featured = ?", *queryParams.IsFeatured)
	}

	if queryParams.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+queryParams.Search+"%", "%"+queryParams.Search+"%")
	}

	// Count total
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply ordering
	switch queryParams.Ordering {
	case "name":
		query = query.Order("name ASC")
	case "price":
		query = query.Order("price ASC")
	case "created_at":
		query = query.Order("created_at DESC")
	case "updated_at":
		query = query.Order("updated_at DESC")
	default:
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	err = query.Limit(queryParams.Limit).Offset(queryParams.Offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, int(totalCount), nil
}

// GetVendorCategories gets categories for a vendor
func (r *VendorProductRepositoryImpl) GetVendorCategories(vendorID uint) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("vendor_id = ? AND is_deleted = ?", vendorID, false).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetOrCreateCategory gets an existing category or creates a new one for the vendor
func (r *VendorProductRepositoryImpl) GetOrCreateCategory(vendorID uint, categoryName string) (*models.Category, error) {
	var category models.Category

	// First, try to find existing category for this vendor
	err := r.db.Where("vendor_id = ? AND name = ? AND is_deleted = ?", vendorID, categoryName, false).First(&category).Error
	if err == nil {
		// Category exists, return it
		return &category, nil
	}

	if err != gorm.ErrRecordNotFound {
		// Some other error occurred
		return nil, err
	}

	// Category doesn't exist, create a new one
	category = models.Category{
		VendorID:  vendorID,
		Name:      categoryName,
		IsDeleted: false,
	}

	err = r.db.Create(&category).Error
	if err != nil {
		return nil, err
	}

	return &category, nil
}
