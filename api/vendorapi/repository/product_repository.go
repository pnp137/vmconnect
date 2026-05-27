package repository

import (
	"gorm.io/gorm"
	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/database/models"
	"linksupply.io/vmconnect/system"
)

type VendorProductRepository interface {
	CreateProduct(vendorID uint, product *models.Product) (*models.Product, error)
	CreateProductVariants(vendorID, productID uint, variants []models.ProductVariant) error
	SyncProductVariants(vendorID, productID uint, variants []models.ProductVariant) error
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

func (r *VendorProductRepositoryImpl) CreateProductVariants(vendorID, productID uint, variants []models.ProductVariant) error {
	var product models.Product
	if err := r.db.Where("id = ? AND vendor_id = ?", productID, vendorID).First(&product).Error; err != nil {
		return err
	}

	for i := range variants {
		variants[i].ProductID = productID
	}
	return r.db.Create(&variants).Error
}

// SyncProductVariants updates/creates variants and deactivates missing existing variants.
// Rules:
// - If variant.ID > 0: update existing variant (must belong to product)
// - If variant.ID == 0: create new variant
// - Any existing DB variant not present in request: set is_active=false
func (r *VendorProductRepositoryImpl) SyncProductVariants(vendorID, productID uint, variants []models.ProductVariant) error {
	var product models.Product
	if err := r.db.Where("id = ? AND vendor_id = ?", productID, vendorID).First(&product).Error; err != nil {
		return err
	}

	var existing []models.ProductVariant
	if err := r.db.Where("product_id = ?", productID).Find(&existing).Error; err != nil {
		return err
	}

	existingByID := make(map[uint]models.ProductVariant, len(existing))
	for _, v := range existing {
		existingByID[v.ID] = v
	}

	seen := map[uint]bool{}

	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range variants {
			v.ProductID = productID
			if v.ID > 0 {
				if _, ok := existingByID[v.ID]; !ok {
					return gorm.ErrRecordNotFound
				}
				seen[v.ID] = true
				if err := tx.Model(&models.ProductVariant{}).
					Where("id = ? AND product_id = ?", v.ID, productID).
					Updates(map[string]interface{}{
						"name":      v.Name,
						"unit":      v.Unit,
						"quantity":  v.Quantity,
						"price":     v.Price,
						"mrp":       v.MRP,
						"moq":       v.MOQ,
						"stock":     v.Stock,
						"is_active": v.IsActive,
					}).Error; err != nil {
					return err
				}
				continue
			}

			if err := tx.Create(&v).Error; err != nil {
				return err
			}
		}

		for _, ex := range existing {
			if !seen[ex.ID] {
				if err := tx.Model(&models.ProductVariant{}).
					Where("id = ? AND product_id = ?", ex.ID, productID).
					Update("is_active", false).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetProductByID gets a product by ID for a specific vendor
func (r *VendorProductRepositoryImpl) GetProductByID(vendorID, productID uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Category").
		Preload("Variants", "is_active = ?", true).
		Where("id = ? AND vendor_id = ?", productID, vendorID).
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
	err := r.db.Where("id = ? AND vendor_id = ?", productID, vendorID).First(&product).Error
	if err != nil {
		return nil, err
	}

	// Update the product
	err = r.db.Model(&models.Product{}).
		Where("id = ? AND vendor_id = ?", productID, vendorID).
		Updates(updates).Error
	if err != nil {
		return nil, err
	}

	// Get the updated product with preloaded category
	err = r.db.Preload("Category").Preload("Variants", "is_active = ?", true).
		Where("id = ? AND vendor_id = ?", productID, vendorID).
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
	err := r.db.Where("id = ? AND vendor_id = ?", productID, vendorID).First(&product).Error
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
	query := r.db.Model(&models.Product{}).
		Preload("Category").
		Preload("Variants", "is_active = ?", true).
		Where("vendor_id = ?", vendorID)

	// Apply filters
	if queryParams.CategoryName != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.vendor_id = ? AND categories.is_deleted = ? AND categories.name = ?", vendorID, false, queryParams.CategoryName)
	}

	if queryParams.ProductID > 0 {
		query = query.Where("products.id = ?", queryParams.ProductID)
	}

	if queryParams.IsActive != nil {
		query = query.Where("is_active = ?", *queryParams.IsActive)
	}

	if queryParams.IsFeatured != nil {
		query = query.Where("is_featured = ?", *queryParams.IsFeatured)
	}

	if queryParams.NameOrDescription != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+queryParams.NameOrDescription+"%", "%"+queryParams.NameOrDescription+"%")
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
		query = query.Joins("LEFT JOIN product_variants pv ON pv.product_id = products.id AND pv.is_active = ?", true).Order("pv.price ASC")
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
