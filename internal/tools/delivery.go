package tools

import (
	"context"
	"fmt"

	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

type DeliveryTool struct {
	DB *gorm.DB
}

type DeliveryStatusResult struct {
	OrderID       uint    `json:"order_id"`
	DeliveryID    uint    `json:"delivery_id"`
	Status        string  `json:"status"`
	ScheduledDate *string `json:"scheduled_date,omitempty"`
	DeliveredAt   *string `json:"delivered_at,omitempty"`
}

func (t *DeliveryTool) GetDeliveryStatus(
	ctx context.Context,
	orderID uint,
) (*DeliveryStatusResult, error) {

	var delivery models.Delivery

	err := t.DB.WithContext(ctx).
		Where("order_id = ?", orderID).
		First(&delivery).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf(
				"delivery information not found for order %d",
				orderID,
			)
		}

		return nil, fmt.Errorf(
			"failed to fetch delivery for order %d: %w",
			orderID,
			err,
		)
	}

	var scheduledDate *string
	var deliveredAt *string

	if delivery.ScheduledDate != nil {
		formatted := delivery.ScheduledDate.Format("2006-01-02T15:04:05Z")
		scheduledDate = &formatted
	}

	if delivery.DeliveredAt != nil {
		formatted := delivery.DeliveredAt.Format("2006-01-02T15:04:05Z")
		deliveredAt = &formatted
	}

	return &DeliveryStatusResult{
		OrderID:       delivery.OrderID,
		DeliveryID:    delivery.ID,
		Status:        delivery.Status,
		ScheduledDate: scheduledDate,
		DeliveredAt:   deliveredAt,
	}, nil
}