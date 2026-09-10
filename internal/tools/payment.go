package tools

import (
	"context"
	"fmt"

	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

type PaymentTool struct {
	DB *gorm.DB
}

type PaymentStatusResult struct {
	OrderID     uint     `json:"order_id"`
	PaymentID   uint     `json:"payment_id"`
	Status      string   `json:"status"`
	Amount      float64  `json:"amount"`
	PaidAt      *string  `json:"paid_at,omitempty"`
}

func (t *PaymentTool) GetPaymentStatus(
	ctx context.Context,
	orderID uint,
) (*PaymentStatusResult, error) {

	var payment models.Payment

	err := t.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&payment).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf(
				"payment information not found for order %d",
				orderID,
			)
		}

		return nil, fmt.Errorf(
			"failed to fetch payment for order %d: %w",
			orderID,
			err,
		)
	}

	var paidAt *string

	if payment.PaidAt != nil {
		formatted := payment.PaidAt.Format("2006-01-02T15:04:05Z")
		paidAt = &formatted
	}

	return &PaymentStatusResult{
		OrderID:   payment.OrderID,
		PaymentID: payment.ID,
		Status:    payment.Status,
		Amount:    payment.Amount,
		PaidAt:    paidAt,
	}, nil
}