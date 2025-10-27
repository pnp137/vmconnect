package models

import "strconv"

type Role int

const (
	ROLE_ADMIN    Role = 1
	ROLE_VENDOR   Role = 2
	ROLE_MERCHANT Role = 3
)

func (role Role) String() string {
	switch role {
	case ROLE_ADMIN:
		return "admin"
	case ROLE_VENDOR:
		return "vendor"
	case ROLE_MERCHANT:
		return "merchant"
	default:
		return "unknown"
	}
}

var (
	roleStringMap = map[string]Role{
		"admin":    ROLE_ADMIN,
		"vendor":   ROLE_VENDOR,
		"merchant": ROLE_MERCHANT,
	}
)

func ParseRoleString(str string) (Role, bool) {
	role, ok := roleStringMap[str]
	return role, ok
}

func (role Role) IsValid() bool {
	return role >= ROLE_ADMIN && role <= ROLE_MERCHANT
}

var RoleMap = map[Role]int{
	ROLE_ADMIN:    1,
	ROLE_VENDOR:   2,
	ROLE_MERCHANT: 3,
}

func CheckIfRoleStringIsValid(str string) bool {
	roleInt, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	for _, v := range RoleMap {
		if v == roleInt {
			return true
		}
	}
	return false
}

