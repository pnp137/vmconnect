package utils

import (
	"fmt"

	"linksupply.io/vmconnect/database/models"
)

// StateTransitionRule defines allowed transitions for a given role
type StateTransitionRule struct {
	FromStatus      models.OrderStatus
	ToStatuses      []models.OrderStatus
	AllowedRoles    []models.ActorRole
	RequiresInvoice bool
	RequiresPayment bool
}

// OrderStateValidator handles order state transition validation
type OrderStateValidator struct {
	transitionRules []StateTransitionRule
}

// NewOrderStateValidator creates a new state validator with predefined rules
func NewOrderStateValidator() *OrderStateValidator {
	validator := &OrderStateValidator{
		transitionRules: []StateTransitionRule{
			// FROM: IN_CART
			{
				FromStatus:   models.ORDER_IN_CART,
				ToStatuses:   []models.OrderStatus{models.ORDER_PLACED, models.ORDER_CANCELLED},
				AllowedRoles: []models.ActorRole{models.ACTOR_MERCHANT, models.ACTOR_SYSTEM},
			},
			// FROM: PLACED
			{
				FromStatus:   models.ORDER_PLACED,
				ToStatuses:   []models.OrderStatus{models.ORDER_PAYMENT_PENDING, models.ORDER_CONFIRMED, models.ORDER_CANCELLED},
				AllowedRoles: []models.ActorRole{models.ACTOR_MERCHANT, models.ACTOR_VENDOR},
			},
			// FROM: PAYMENT_PENDING
			{
				FromStatus:   models.ORDER_PAYMENT_PENDING,
				ToStatuses:   []models.OrderStatus{models.ORDER_CONFIRMED, models.ORDER_CANCELLED},
				AllowedRoles: []models.ActorRole{models.ACTOR_VENDOR, models.ACTOR_MERCHANT},
			},
			// FROM: CONFIRMED
			{
				FromStatus:   models.ORDER_CONFIRMED,
				ToStatuses:   []models.OrderStatus{models.ORDER_INVOICED, models.ORDER_SHIPPED, models.ORDER_CANCELLED},
				AllowedRoles: []models.ActorRole{models.ACTOR_VENDOR},
			},
			// FROM: INVOICED
			{
				FromStatus:      models.ORDER_INVOICED,
				ToStatuses:      []models.OrderStatus{models.ORDER_SHIPPED, models.ORDER_CANCELLED},
				AllowedRoles:    []models.ActorRole{models.ACTOR_VENDOR},
				RequiresInvoice: true,
			},
			// FROM: SHIPPED
			{
				FromStatus:   models.ORDER_SHIPPED,
				ToStatuses:   []models.OrderStatus{models.ORDER_DELIVERED},
				AllowedRoles: []models.ActorRole{models.ACTOR_MERCHANT, models.ACTOR_VENDOR},
			},
			// FROM: DELIVERED
			{
				FromStatus:   models.ORDER_DELIVERED,
				ToStatuses:   []models.OrderStatus{models.ORDER_COMPLETED},
				AllowedRoles: []models.ActorRole{models.ACTOR_MERCHANT, models.ACTOR_VENDOR},
			},
			// FROM: ON_HOLD
			{
				FromStatus:   models.ORDER_ON_HOLD,
				ToStatuses:   []models.OrderStatus{models.ORDER_CONFIRMED, models.ORDER_CANCELLED},
				AllowedRoles: []models.ActorRole{models.ACTOR_VENDOR, models.ACTOR_SYSTEM},
			},
		},
	}
	return validator
}

// ValidateTransition checks if a state transition is allowed
func (v *OrderStateValidator) ValidateTransition(
	currentStatus models.OrderStatus,
	targetStatus models.OrderStatus,
	actorRole models.ActorRole,
) error {
	// Find applicable rule
	var applicableRule *StateTransitionRule
	for _, rule := range v.transitionRules {
		if rule.FromStatus == currentStatus {
			applicableRule = &rule
			break
		}
	}

	if applicableRule == nil {
		return fmt.Errorf("no transition rules defined for status %s", currentStatus.String())
	}

	// Check if target status is allowed
	isTargetAllowed := false
	for _, allowedStatus := range applicableRule.ToStatuses {
		if allowedStatus == targetStatus {
			isTargetAllowed = true
			break
		}
	}

	if !isTargetAllowed {
		return fmt.Errorf(
			"cannot transition from %s to %s",
			currentStatus.String(),
			targetStatus.String(),
		)
	}

	// Check if role is allowed
	isRoleAllowed := false
	for _, allowedRole := range applicableRule.AllowedRoles {
		if allowedRole == actorRole {
			isRoleAllowed = true
			break
		}
	}

	if !isRoleAllowed {
		return fmt.Errorf(
			"role %s is not allowed to transition from %s to %s",
			actorRole.String(),
			currentStatus.String(),
			targetStatus.String(),
		)
	}

	return nil
}

// GetAllowedTransitions returns all allowed target statuses from current status for a role
func (v *OrderStateValidator) GetAllowedTransitions(
	currentStatus models.OrderStatus,
	actorRole models.ActorRole,
) []models.OrderStatus {
	var allowed []models.OrderStatus

	for _, rule := range v.transitionRules {
		if rule.FromStatus == currentStatus {
			// Check if role is allowed
			for _, allowedRole := range rule.AllowedRoles {
				if allowedRole == actorRole {
					allowed = append(allowed, rule.ToStatuses...)
					break
				}
			}
			break
		}
	}

	return allowed
}

// CanRoleUpdateStatus checks if a role can update to a specific status
func (v *OrderStateValidator) CanRoleUpdateStatus(
	targetStatus models.OrderStatus,
	actorRole models.ActorRole,
) bool {
	// Certain statuses are exclusive to certain roles
	switch targetStatus {
	case models.ORDER_IN_CART:
		return actorRole == models.ACTOR_MERCHANT || actorRole == models.ACTOR_SYSTEM

	case models.ORDER_PLACED, models.ORDER_PAYMENT_PENDING, models.ORDER_DELIVERED, models.ORDER_COMPLETED:
		return actorRole == models.ACTOR_MERCHANT || actorRole == models.ACTOR_SYSTEM

	case models.ORDER_CONFIRMED, models.ORDER_INVOICED, models.ORDER_SHIPPED:
		return actorRole == models.ACTOR_VENDOR || actorRole == models.ACTOR_SYSTEM

	case models.ORDER_CANCELLED, models.ORDER_ON_HOLD:
		// Both can cancel or put on hold
		return actorRole == models.ACTOR_MERCHANT || actorRole == models.ACTOR_VENDOR || actorRole == models.ACTOR_SYSTEM

	default:
		return actorRole == models.ACTOR_SYSTEM
	}
}
