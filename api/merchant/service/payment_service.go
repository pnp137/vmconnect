package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"linksupply.io/vmconnect/api/merchant/dto"
	"linksupply.io/vmconnect/api/merchant/repository"
	"linksupply.io/vmconnect/database/models"
)

// PaymentService handles payment-related operations for merchants
type PaymentService interface {
	SubmitPayment(ctx context.Context, merchantID uint, orderID uint, req *dto.SubmitPaymentRequest) (*dto.SubmitPaymentResponse, error)
	GetOrderPayments(ctx context.Context, merchantID uint, orderID uint) (*dto.GetPaymentsResponse, error)
}

type PaymentServiceImpl struct {
	repo repository.MerchantRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(repo repository.MerchantRepository) PaymentService {
	return &PaymentServiceImpl{repo: repo}
}

// SubmitPayment allows merchant to submit a payment for an order
func (s *PaymentServiceImpl) SubmitPayment(
	ctx context.Context,
	merchantID uint,
	orderID uint,
	req *dto.SubmitPaymentRequest,
) (*dto.SubmitPaymentResponse, error) {
	// 1. Get merchant details
	merchant, err := s.repo.GetMerchantByID(merchantID)
	if err != nil {
		return nil, errors.New("merchant not found")
	}

	// 2. Get order and verify ownership
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if !orderBelongsToMerchant(order, merchantID) {
		return nil, errors.New("unauthorized: order does not belong to this merchant")
	}

	// 3. Validate and parse payment mode enum
	paymentMode := models.PaymentMode(req.PaymentMode)
	if !models.CheckIfPaymentModeStringIsValid(strconv.Itoa(int(paymentMode))) {
		return nil, errors.New("invalid payment mode")
	}

	// 4. Create payment record
	actorRole := models.ACTOR_MERCHANT
	paymentStatus := models.PaymentStatusPending
	now := time.Now()

	payment := &models.Payment{
		OrderID:         orderID,
		SubmittedBy:     merchant.OwnerID,
		ActorRole:       &actorRole,
		Amount:          req.Amount,
		PaymentMode:     &paymentMode,
		PaymentType:     req.PaymentType,
		UTR:             req.UTR,
		ReferenceNumber: req.ReferenceNumber,
		TransactionID:   req.TransactionID,
		Status:          &paymentStatus,
		Verified:        false,
		Notes:           req.Notes,
		Metadata:        req.Metadata,
		PaidAt:          now,
	}

	createdPayment, err := s.repo.CreatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	// 5. Auto-recalculate order amounts
	err = s.repo.RecalculateOrderPayments(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to recalculate order payments: %w", err)
	}

	// 6. Get updated order
	updatedOrder, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, err
	}

	// 7. Update order status to PAYMENT_PENDING if currently PLACED
	if *updatedOrder.Status == models.ORDER_PLACED {
		paymentPending := models.ORDER_PAYMENT_PENDING
		updatedOrder.Status = &paymentPending
		updatedOrder.StatusUpdatedAt = &now
		updatedOrder.StatusUpdatedBy = merchant.OwnerID
		err = s.repo.UpdateOrderStatus(orderID, paymentPending)
		if err != nil {
			return nil, fmt.Errorf("failed to update order status: %w", err)
		}
	}

	// 8. Create order activity log
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    merchant.OwnerID,
		ActorRole:  &actorRole,
		OrderState: updatedOrder.Status,
		Remarks:    fmt.Sprintf("Payment of ₹%.2f submitted via %s", req.Amount, paymentMode.String()),
		Metadata: models.JSONMap{
			"action":           "submit_payment",
			"payment_id":       createdPayment.ID,
			"amount":           req.Amount,
			"payment_mode":     paymentMode.String(),
			"payment_mode_int": int(paymentMode),
			"payment_type":     req.PaymentType,
			"total_paid":       updatedOrder.PaidAmount,
			"outstanding":      updatedOrder.OutstandingAmount,
		},
	}

	_, err = s.repo.CreateOrderActivity(activity)
	if err != nil {
		// Log error but don't fail the payment submission
		fmt.Printf("Warning: failed to create order activity: %v\n", err)
	}

	// 9. Prepare response
	return &dto.SubmitPaymentResponse{
		PaymentID:     createdPayment.ID,
		OrderID:       orderID,
		Amount:        createdPayment.Amount,
		PaymentMode:   paymentMode.String(),
		PaymentStatus: paymentStatus.String(),
		TotalPaid:     updatedOrder.PaidAmount,
		Outstanding:   updatedOrder.OutstandingAmount,
		OrderStatus:   updatedOrder.Status.String(),
		OrderPaid:     updatedOrder.OrderPaid,
		SubmittedAt:   createdPayment.PaidAt,
		Message:       fmt.Sprintf("Payment recorded successfully. Balance: ₹%.2f", updatedOrder.OutstandingAmount),
	}, nil
}

// GetOrderPayments retrieves all payments for an order
func (s *PaymentServiceImpl) GetOrderPayments(
	ctx context.Context,
	merchantID uint,
	orderID uint,
) (*dto.GetPaymentsResponse, error) {
	// 1. Get order and verify ownership
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if !orderBelongsToMerchant(order, merchantID) {
		return nil, errors.New("unauthorized: order does not belong to this merchant")
	}

	// 2. Get all payments for the order
	payments, err := s.repo.GetOrderPayments(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}

	// 3. Build payment summaries
	paymentSummaries := make([]dto.PaymentSummary, 0, len(payments))
	byMode := make(map[string]float64)
	var verified, pendingVerification float64

	for _, p := range payments {
		summary := dto.PaymentSummary{
			ID:              p.ID,
			Amount:          p.Amount,
			PaymentMode:     p.PaymentMode.String(),
			PaymentType:     p.PaymentType,
			PaymentStatus:   p.Status.String(),
			Verified:        p.Verified,
			UTR:             p.UTR,
			ReferenceNumber: p.ReferenceNumber,
			TransactionID:   p.TransactionID,
			Notes:           p.Notes,
			PaidAt:          p.PaidAt,
			VerifiedAt:      p.VerifiedAt,
		}
		paymentSummaries = append(paymentSummaries, summary)

		// Aggregate by mode
		mode := p.PaymentMode.String()
		byMode[mode] += p.Amount

		// Calculate verified vs pending
		if p.Verified {
			verified += p.Amount
		} else if *p.Status != models.PaymentStatusFailed && *p.Status != models.PaymentStatusCancelled {
			pendingVerification += p.Amount
		}
	}

	// 4. Build response
	return &dto.GetPaymentsResponse{
		OrderID:           orderID,
		OrderStatus:       order.Status.String(),
		TotalAmount:       order.TotalAmount,
		PaidAmount:        order.PaidAmount,
		OutstandingAmount: order.OutstandingAmount,
		OrderPaid:         order.OrderPaid,
		OrderPaidAt:       order.OrderPaidAt,
		Payments:          paymentSummaries,
		PaymentSummary: dto.PaymentModeSummary{
			ByMode:              byMode,
			Verified:            verified,
			PendingVerification: pendingVerification,
		},
	}, nil
}
