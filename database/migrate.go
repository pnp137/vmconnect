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
		&models.OrderActivity{},

		// Payment and invoice entities
		&models.Payment{},
		&models.Invoice{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrations completed successfully")
}
