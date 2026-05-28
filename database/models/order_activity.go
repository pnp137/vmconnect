package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type OrderActivity struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	OrderID    uint         `gorm:"not null;index:idx_order_activity_state_history,priority:1" json:"orderId"`
	ActorID    uint         `gorm:"not null;comment:'User ID from users table'" json:"actorId"`
	ActorRole  *ActorRole   `gorm:"type:int;not null" json:"actorRole"`                                   // merchant, vendor, system
	OrderState *OrderStatus `gorm:"index:idx_order_activity_state_history,priority:2" json:"order_state"` // Order status enum
	Remarks    string       `gorm:"type:text" json:"remarks,omitempty"`
	Metadata   JSONMap      `gorm:"type:json" json:"metadata,omitempty"`
	CreatedAt  time.Time    `gorm:"index:idx_order_activity_state_history,priority:4" json:"createdAt"`
	IsDeleted  bool         `gorm:"default:false;index:idx_order_activity_state_history,priority:3" json:"isDeleted"`
	UpdatedAt  time.Time    `json:"updatedAt"`

	// Relationships
	Actor User `gorm:"foreignKey:ActorID" json:"actor,omitempty"`
}

// JSONMap is a custom type for JSON fields
type JSONMap map[string]interface{}

// Scan implements sql.Scanner interface
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value implements driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
