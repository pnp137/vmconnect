package database

import (
	"log"

	"gorm.io/gorm"
	"linksupply.io/vmconnect/database/models"
)

// Migrate performs database migrations
func Migrate(db *gorm.DB) {
	log.Println("Running database migrations...")

	// Migrate all models in the correct order (respecting foreign key dependencies)
	err := db.AutoMigrate(
		// Core entities first (no foreign keys)
		&models.Role{},
		&models.User{},

		// Entities that depend on User
		&models.Vendor{},
		&models.Merchant{},

		// Entities that depend on Vendor
		&models.Category{},
		&models.Product{},

		// Junction/relationship tables
		&models.VendorMerchantMapping{},
		&models.MerchantProductVisibility{},

		// Order related entities
		&models.Order{},
		&models.OrderItem{},

		// Payment and invoice entities
		&models.Payment{},
		&models.Invoice{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed initial data
	seedInitialData(db)

	log.Println("Database migrations completed successfully")
}

// seedInitialData creates initial required data
func seedInitialData(db *gorm.DB) {
	log.Println("Seeding initial data...")

	// Create default roles if they don't exist
	roles := []models.Role{
		{Name: "admin", Description: "System Administrator"},
		{Name: "vendor", Description: "Product Vendor/Supplier"},
		{Name: "merchant", Description: "Product Merchant/Retailer"},
	}

	for _, role := range roles {
		var existingRole models.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("Created role: %s", role.Name)
			}
		}
	}

	log.Println("Initial data seeding completed")
}
