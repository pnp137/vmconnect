package models

import (
	"database/sql/driver"
	"errors"
)

type ActorRole int

const (
	ACTOR_MERCHANT ActorRole = 1
	ACTOR_VENDOR   ActorRole = 2
	ACTOR_SYSTEM   ActorRole = 3
)

// String returns the string representation of the ActorRole
func (ar ActorRole) String() string {
	switch ar {
	case ACTOR_MERCHANT:
		return "merchant"
	case ACTOR_VENDOR:
		return "vendor"
	case ACTOR_SYSTEM:
		return "system"
	default:
		return "unknown"
	}
}

// MarshalJSON converts ActorRole to JSON string
func (ar ActorRole) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ar.String() + `"`), nil
}

// UnmarshalJSON converts JSON string to ActorRole
func (ar *ActorRole) UnmarshalJSON(data []byte) error {
	str := string(data)
	// Remove quotes if present
	if len(str) > 0 && str[0] == '"' {
		str = str[1 : len(str)-1]
	}

	parsedRole, isValid := ParseActorRoleString(str)
	if !isValid {
		return errors.New("invalid actor role")
	}

	*ar = parsedRole
	return nil
}

// Scan implements the sql.Scanner interface
func (ar *ActorRole) Scan(value interface{}) error {
	if value == nil {
		return errors.New("actor role cannot be null")
	}

	switch v := value.(type) {
	case int64:
		*ar = ActorRole(v)
		return nil
	case []byte:
		str := string(v)
		parsedRole, isValid := ParseActorRoleString(str)
		if !isValid {
			return errors.New("invalid actor role")
		}
		*ar = parsedRole
		return nil
	case string:
		parsedRole, isValid := ParseActorRoleString(v)
		if !isValid {
			return errors.New("invalid actor role")
		}
		*ar = parsedRole
		return nil
	default:
		return errors.New("invalid actor role type")
	}
}

// Value implements the driver.Valuer interface
func (ar ActorRole) Value() (driver.Value, error) {
	return int64(ar), nil
}

// ParseActorRoleString converts string to ActorRole
func ParseActorRoleString(s string) (ActorRole, bool) {
	switch s {
	case "merchant":
		return ACTOR_MERCHANT, true
	case "vendor":
		return ACTOR_VENDOR, true
	case "system":
		return ACTOR_SYSTEM, true
	default:
		return 0, false
	}
}

// IsValid checks if the ActorRole is valid
func (ar ActorRole) IsValid() bool {
	switch ar {
	case ACTOR_MERCHANT, ACTOR_VENDOR, ACTOR_SYSTEM:
		return true
	default:
		return false
	}
}
