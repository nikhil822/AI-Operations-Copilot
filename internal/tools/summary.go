package tools

import (
	"context"
	"fmt"

	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

type SummaryTool struct {
	DB *gorm.DB
}

type FullOrderSummary struct {
	Order struct {
		ID        uint   `json:"id"`
		Status    string `json:"status"`
		OrderDate string `json:"order_date"`
	} `json:"order"`

	Customer struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	} `json:"customer"`

	Vehicle struct {
		ID    uint    `json:"id"`
		Make  string  `json:"make"`
		Model string  `json:"model"`
		Year  int     `json:"year"`
		Price float64 `json:"price"`
	} `json:"vehicle"`

	Payment *PaymentStatusResult `json:"payment,omitempty"`

	Delivery *DeliveryStatusResult `json:"delivery,omitempty"`
}

func (t *SummaryTool) GetFullOrderSummary(
	ctx context.Context,
	orderID uint,
) (*FullOrderSummary, error) {

	var order models.Order

	// get the order along with its associated customer, vehicle, payment, and delivery information
	err := t.DB.WithContext(ctx).
		Preload("Customer").
		Preload("Vehicle").
		Preload("Payment").
		Preload("Delivery").
		First(&order, orderID).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("order %d not found", orderID)
		}

		return nil, fmt.Errorf(
			"failed to fetch full order summary for %d: %w",
			orderID,
			err,
		)
	}

	result := &FullOrderSummary{}

	result.Order.ID = order.ID
	result.Order.Status = order.Status
	result.Order.OrderDate = order.OrderDate.Format("2006-01-02")

	result.Customer.ID = order.Customer.ID
	result.Customer.Name = order.Customer.Name
	result.Customer.Phone = order.Customer.Phone
	result.Customer.Email = order.Customer.Email

	result.Vehicle.ID = order.Vehicle.ID
	result.Vehicle.Make = order.Vehicle.Make
	result.Vehicle.Model = order.Vehicle.Model
	result.Vehicle.Year = order.Vehicle.Year
	result.Vehicle.Price = order.Vehicle.Price

	if order.Payment != nil {
		var paidAt *string

		if order.Payment.PaidAt != nil {
			formatted := order.Payment.PaidAt.Format(
				"2006-01-02T15:04:05Z",
			)
			paidAt = &formatted
		}

		result.Payment = &PaymentStatusResult{
			OrderID:   order.Payment.OrderID,
			PaymentID: order.Payment.ID,
			Status:    order.Payment.Status,
			Amount:    order.Payment.Amount,
			PaidAt:    paidAt,
		}
	}

	if order.Delivery != nil {
		var scheduledDate *string
		var deliveredAt *string

		if order.Delivery.ScheduledDate != nil {
			formatted := order.Delivery.ScheduledDate.Format(
				"2006-01-02T15:04:05Z",
			)
			scheduledDate = &formatted
		}

		if order.Delivery.DeliveredAt != nil {
			formatted := order.Delivery.DeliveredAt.Format(
				"2006-01-02T15:04:05Z",
			)
			deliveredAt = &formatted
		}

		result.Delivery = &DeliveryStatusResult{
			OrderID:       order.Delivery.OrderID,
			DeliveryID:    order.Delivery.ID,
			Status:        order.Delivery.Status,
			ScheduledDate: scheduledDate,
			DeliveredAt:   deliveredAt,
		}
	}

	return result, nil
}