package models

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"size:100;uniqueIndex" json:"email"`
	Phone        string    `gorm:"size:15;uniqueIndex" json:"phone"`
	PasswordHash string    `gorm:"size:255" json:"password_hash"`
	RoleID       *Role     `gorm:"type:int;not null" json:"role_id"` // Role enum: 1=admin, 2=vendor, 3=merchant
	IsDeleted    bool      `gorm:"default:false" json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	IsActive     bool      `gorm:"default:false" json:"is_active"`
}
