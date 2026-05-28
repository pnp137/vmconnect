package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"linksupply.io/vmconnect/api/vendorapi/dto"
	"linksupply.io/vmconnect/api/vendorapi/repository"
	"linksupply.io/vmconnect/database/models"
)

// PaymentService handles payment-related operations for vendors
type PaymentService interface {
	MarkOrderAsPaid(ctx context.Context, vendorID uint, orderID uint, req *dto.MarkOrderPaidRequest) (*dto.MarkOrderPaidResponse, error)
	VerifyPayment(ctx context.Context, vendorID uint, orderID uint, paymentID uint64, req *dto.VerifyPaymentRequest) (*dto.VerifyPaymentResponse, error)
	GetOrderPayments(ctx context.Context, vendorID uint, orderID uint) (*dto.GetOrderPaymentsResponse, error)
}

type PaymentServiceImpl struct {
	repo repository.VendorRepository
}

// NewPaymentService creates a new payment service instance
func NewPaymentService(repo repository.VendorRepository) PaymentService {
	return &PaymentServiceImpl{repo: repo}
}

// MarkOrderAsPaid marks all payments of an order as verified (bulk action)
func (s *PaymentServiceImpl) MarkOrderAsPaid(
	ctx context.Context,
	vendorID uint,
	orderID uint,
	req *dto.MarkOrderPaidRequest,
) (*dto.MarkOrderPaidResponse, error) {
	// 1. Get vendor details
	vendor, err := s.repo.GetVendorByID(vendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	// 2. Get order and verify ownership
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.VendorID != vendorID {
		return nil, errors.New("unauthorized: order does not belong to this vendor")
	}

	// 3. Mark order as paid
	now := time.Now()
	order.OrderPaid = true
	order.OrderPaidAt = &now
	order.OrderPaidBy = vendor.UserID

	// 4. Get all payments and mark them as verified
	payments, err := s.repo.GetOrderPayments(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}

	verifiedCount := 0
	completedStatus := models.PaymentStatusCompleted

	for i := range payments {
		if !payments[i].Verified && *payments[i].Status == models.PaymentStatusPending {
			payments[i].Verified = true
			payments[i].VerifiedBy = vendor.UserID
			payments[i].VerifiedAt = &now
			payments[i].Status = &completedStatus

			err = s.repo.UpdatePayment(&payments[i])
			if err != nil {
				return nil, fmt.Errorf("failed to update payment %d: %w", payments[i].ID, err)
			}
			verifiedCount++
		}
	}

	// 5. Update order status if PAYMENT_PENDING
	if *order.Status == models.ORDER_PAYMENT_PENDING {
		confirmedStatus := models.ORDER_CONFIRMED
		order.Status = &confirmedStatus
		order.StatusUpdatedAt = &now
		order.StatusUpdatedBy = vendor.UserID
	}

	err = s.repo.UpdateOrder(order)
	if err != nil {
		return nil, fmt.Errorf("failed to update order: %w", err)
	}

	// 6. Create order activity log
	actorRole := models.ACTOR_VENDOR
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: order.Status,
		Remarks:    fmt.Sprintf("Order marked as paid - Total: ₹%.2f (%d payments verified)", order.PaidAmount, verifiedCount),
		Metadata: models.JSONMap{
			"action":         "mark_order_paid",
			"payments_count": len(payments),
			"verified_count": verifiedCount,
			"total_amount":   order.PaidAmount,
			"notes":          req.Notes,
			"marked_at":      now,
		},
	}

	_, err = s.repo.CreateOrderActivity(activity)
	if err != nil {
		fmt.Printf("Warning: failed to create order activity: %v\n", err)
	}

	// 7. Prepare response
	return &dto.MarkOrderPaidResponse{
		OrderID:          orderID,
		OrderPaid:        true,
		PaymentsVerified: verifiedCount,
		TotalAmount:      order.PaidAmount,
		MarkedAt:         now,
		OrderStatus:      order.Status.String(),
		Message:          fmt.Sprintf("Order marked as paid successfully. %d payments verified.", verifiedCount),
	}, nil
}

// VerifyPayment verifies an individual payment (optional - rarely used)
func (s *PaymentServiceImpl) VerifyPayment(
	ctx context.Context,
	vendorID uint,
	orderID uint,
	paymentID uint64,
	req *dto.VerifyPaymentRequest,
) (*dto.VerifyPaymentResponse, error) {
	// 1. Get vendor details
	vendor, err := s.repo.GetVendorByID(vendorID)
	if err != nil {
		return nil, errors.New("vendor not found")
	}

	// 2. Get order and verify ownership
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.VendorID != vendorID {
		return nil, errors.New("unauthorized: order does not belong to this vendor")
	}

	// 3. Get payment
	payment, err := s.repo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.OrderID != orderID {
		return nil, errors.New("payment does not belong to this order")
	}

	// 4. Update payment status
	now := time.Now()
	var newStatus models.PaymentStatus

	if req.Verified {
		newStatus = models.PaymentStatusCompleted
		payment.Verified = true
		payment.VerifiedBy = vendor.UserID
		payment.VerifiedAt = &now
	} else {
		newStatus = models.PaymentStatusFailed
		payment.Verified = false
	}

	payment.Status = &newStatus
	if req.Notes != "" {
		payment.Notes = req.Notes
	}

	err = s.repo.UpdatePayment(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to update payment: %w", err)
	}

	// 5. Recalculate order payments (in case failed payment was excluded)
	err = s.repo.RecalculateOrderPayments(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to recalculate order payments: %w", err)
	}

	// 6. Create order activity log
	actorRole := models.ACTOR_VENDOR
	activity := &models.OrderActivity{
		OrderID:    orderID,
		ActorID:    vendor.UserID,
		ActorRole:  &actorRole,
		OrderState: order.Status,
		Remarks:    fmt.Sprintf("Payment #%d %s - ₹%.2f", paymentID, newStatus.String(), payment.Amount),
		Metadata: models.JSONMap{
			"action":         "verify_payment",
			"payment_id":     paymentID,
			"amount":         payment.Amount,
			"payment_mode":   payment.PaymentMode.String(),
			"verified":       req.Verified,
			"payment_status": newStatus.String(),
			"notes":          req.Notes,
		},
	}

	_, err = s.repo.CreateOrderActivity(activity)
	if err != nil {
		fmt.Printf("Warning: failed to create order activity: %v\n", err)
	}

	// 7. Prepare response
	return &dto.VerifyPaymentResponse{
		PaymentID:     paymentID,
		OrderID:       orderID,
		Amount:        payment.Amount,
		PaymentStatus: newStatus.String(),
		Verified:      payment.Verified,
		VerifiedAt:    now,
		Message:       fmt.Sprintf("Payment %s successfully", newStatus.String()),
	}, nil
}

// GetOrderPayments retrieves all payments for an order (vendor view)
func (s *PaymentServiceImpl) GetOrderPayments(
	ctx context.Context,
	vendorID uint,
	orderID uint,
) (*dto.GetOrderPaymentsResponse, error) {
	// 1. Get order and verify ownership
	order, err := s.repo.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.VendorID != vendorID {
		return nil, errors.New("unauthorized: order does not belong to this vendor")
	}

	// 2. Get all payments for the order
	payments, err := s.repo.GetOrderPayments(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payments: %w", err)
	}

	// 3. Build payment details
	paymentDetails := make([]dto.PaymentDetail, 0, len(payments))
	var totalSubmitted, verified, pending float64

	for _, p := range payments {
		detail := dto.PaymentDetail{
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
			SubmittedAt:     p.PaidAt,
			VerifiedAt:      p.VerifiedAt,
		}
		paymentDetails = append(paymentDetails, detail)

		// Calculate summary
		totalSubmitted += p.Amount
		if p.Verified {
			verified += p.Amount
		} else if *p.Status != models.PaymentStatusFailed && *p.Status != models.PaymentStatusCancelled {
			pending += p.Amount
		}
	}

	// 4. Build response
	merchantName := ""
	if order.Merchant.ID > 0 {
		merchantName = order.Merchant.ShopName
	}

	return &dto.GetOrderPaymentsResponse{
		OrderID:           orderID,
		MerchantName:      merchantName,
		OrderStatus:       order.Status.String(),
		TotalAmount:       order.TotalAmount,
		PaidAmount:        order.PaidAmount,
		OutstandingAmount: order.OutstandingAmount,
		OrderPaid:         order.OrderPaid,
		OrderPaidAt:       order.OrderPaidAt,
		Payments:          paymentDetails,
		Summary: dto.PaymentStatusSummary{
			TotalSubmitted:      totalSubmitted,
			Verified:            verified,
			PendingVerification: pending,
			PaymentsCount:       len(payments),
		},
	}, nil
}
