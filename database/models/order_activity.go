package models

import (
	"time"
)

type OrderActivity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `json:"orderId"`
	ActorType string    `gorm:"size:50" json:"actorType"` // merchant, vendor, system
	ActorID   uint      `json:"actorId"`
	Remarks   string    `gorm:"type:text" json:"remarks,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	IsDeleted bool      `gorm:"default:false" json:"isDeleted"`
	UpdatedAt time.Time `json:"updatedAt"`
}
