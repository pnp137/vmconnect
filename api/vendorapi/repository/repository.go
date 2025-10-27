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
	GetVendorByID(vendorID uint) (*models.Vendor, error)
	// Order operations
	GetOrderByID(orderID uint) (*models.Order, error)
	UpdateOrderStatus(orderID uint, status models.OrderStatus) error
	CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error)
	CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error)
	// Payment operations
	CreatePayment(payment *models.Payment) (*models.Payment, error)
	UpdatePayment(payment *models.Payment) error
	GetPaymentByID(paymentID uint64) (*models.Payment, error)
	GetOrderPayments(orderID uint) ([]models.Payment, error)
	RecalculateOrderPayments(orderID uint) error
	UpdateOrder(order *models.Order) error
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

// GetVendorByID gets vendor by ID with user info preloaded
func (r *VendorRepositoryImpl) GetVendorByID(vendorID uint) (*models.Vendor, error) {
	var vendor models.Vendor
	err := r.db.Preload("User").Where("id = ? AND is_deleted = ?", vendorID, false).First(&vendor).Error
	if err != nil {
		return nil, err
	}
	return &vendor, nil
}

// GetOrderByID gets order by ID with related data
func (r *VendorRepositoryImpl) GetOrderByID(orderID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Model(&models.Order{}).Preload("Vendor").Preload("Merchant").Preload("OrderItems.Product").
		Where("id = ? AND is_deleted = ?", orderID, false).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// UpdateOrderStatus updates order status
func (r *VendorRepositoryImpl) UpdateOrderStatus(orderID uint, status models.OrderStatus) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// CreateOrderActivity creates order activity record
func (r *VendorRepositoryImpl) CreateOrderActivity(activity *models.OrderActivity) (*models.OrderActivity, error) {
	err := r.db.Create(activity).Error
	if err != nil {
		return nil, err
	}
	return activity, nil
}

// CheckStatusInHistory checks if a status was previously recorded in order activity
func (r *VendorRepositoryImpl) CheckStatusInHistory(orderID uint, status models.OrderStatus) (bool, *models.OrderActivity, error) {
	var activity models.OrderActivity

	err := r.db.Where(
		"order_id = ? AND order_state = ? AND is_deleted = ?",
		orderID,
		int(status),
		false,
	).Order("created_at DESC").First(&activity).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil, nil // Not found in history
		}
		return false, nil, err // Database error
	}

	return true, &activity, nil
}

// CreatePayment creates a new payment record
func (r *VendorRepositoryImpl) CreatePayment(payment *models.Payment) (*models.Payment, error) {
	if err := r.db.Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

// UpdatePayment updates an existing payment record
func (r *VendorRepositoryImpl) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// GetPaymentByID retrieves a payment by ID
func (r *VendorRepositoryImpl) GetPaymentByID(paymentID uint64) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ? AND is_deleted = ?", paymentID, false).
		Preload("Order").
		Preload("SubmittedByUser").
		First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetOrderPayments retrieves all payments for an order
func (r *VendorRepositoryImpl) GetOrderPayments(orderID uint) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Where("order_id = ? AND is_deleted = ?", orderID, false).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, err
	}
	return payments, nil
}

// RecalculateOrderPayments recalculates total paid and outstanding amounts for an order
func (r *VendorRepositoryImpl) RecalculateOrderPayments(orderID uint) error {
	// Get all non-failed, non-cancelled payments
	var payments []models.Payment
	err := r.db.Where("order_id = ? AND is_deleted = ? AND status NOT IN (?)",
		orderID, false, []int{int(models.PaymentStatusFailed), int(models.PaymentStatusCancelled)}).
		Find(&payments).Error
	if err != nil {
		return err
	}

	// Calculate total paid
	totalPaid := 0.0
	for _, payment := range payments {
		totalPaid += payment.Amount
	}

	// Update order
	var order models.Order
	err = r.db.First(&order, orderID).Error
	if err != nil {
		return err
	}

	order.PaidAmount = totalPaid
	order.OutstandingAmount = order.TotalAmount - totalPaid

	return r.db.Save(&order).Error
}

// UpdateOrder updates an order record
func (r *VendorRepositoryImpl) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}
