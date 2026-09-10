package tools

import (
	"context"
	"fmt"

	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

type OrderTool struct {
	DB *gorm.DB
}

type OrderStatusResult struct {
	OrderID    uint   `json:"order_id"`
	Status     string `json:"status"`
	OrderDate  string `json:"order_date"`
	CustomerID uint   `json:"customer_id"`
	VehicleID  uint   `json:"vehicle_id"`
}

func (t *OrderTool) GetOrderStatus(
	ctx context.Context,
	orderID uint,
) (*OrderStatusResult, error) {

	var order models.Order

	err := t.DB.WithContext(ctx).
		First(&order, orderID).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order %d not found", orderID)
		}

		return nil, fmt.Errorf("failed to fetch order %d: %w", orderID, err)
	}

	return &OrderStatusResult{
		OrderID:    order.ID,
		Status:     order.Status,
		OrderDate:  order.OrderDate.Format("2006-01-02"),
		CustomerID: order.CustomerID,
		VehicleID:  order.VehicleID,
	}, nil
}