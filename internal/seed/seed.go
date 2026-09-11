package seed

import (
	"fmt"
	"time"

	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	if err := clearDatabase(db); err != nil {
		return err
	}

	customers := make([]models.Customer, 10)
	for index := range customers {
		id := uint(index + 1)
		customers[index] = models.Customer{
			ID:    id,
			Name:  fmt.Sprintf("Customer %d", id),
			Phone: fmt.Sprintf("+91-98765432%02d", 10+index),
			Email: fmt.Sprintf("customer%d@example.com", id),
		}
	}
	if err := db.Create(&customers).Error; err != nil {
		return fmt.Errorf("creating customers: %w", err)
	}

	vehicles := []models.Vehicle{
		{ID: 1, Make: "Maruti Suzuki", Model: "Swift", Year: 2023, Price: 750000},
		{ID: 2, Make: "Hyundai", Model: "Creta", Year: 2024, Price: 1450000},
		{ID: 3, Make: "Tata", Model: "Nexon", Year: 2023, Price: 1050000},
		{ID: 4, Make: "Mahindra", Model: "XUV700", Year: 2024, Price: 1850000},
		{ID: 5, Make: "Toyota", Model: "Fortuner", Year: 2022, Price: 3800000},
		{ID: 6, Make: "Honda", Model: "City", Year: 2023, Price: 1350000},
		{ID: 7, Make: "Kia", Model: "Seltos", Year: 2024, Price: 1600000},
		{ID: 8, Make: "Volkswagen", Model: "Taigun", Year: 2023, Price: 1550000},
		{ID: 9, Make: "Skoda", Model: "Slavia", Year: 2024, Price: 1700000},
		{ID: 10, Make: "MG", Model: "Hector", Year: 2023, Price: 1950000},
	}
	if err := db.Create(&vehicles).Error; err != nil {
		return fmt.Errorf("creating vehicles: %w", err)
	}

	orderIDs := []uint{1289, 2231, 4521, 3345, 5678, 6789, 7890, 8901, 9012, 9123, 9345, 9456, 9567, 9678, 9789}
	customerIDs := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 4, 6, 8}
	vehicleIDs := []uint{2, 5, 3, 4, 7, 1, 6, 8, 9, 10, 1, 3, 6, 8, 7}
	statuses := []string{
		models.OrderStatusConfirmed, models.OrderStatusCompleted, models.OrderStatusConfirmed,
		models.OrderStatusPending, models.OrderStatusConfirmed, models.OrderStatusConfirmed,
		models.OrderStatusPending, models.OrderStatusConfirmed, models.OrderStatusCompleted,
		models.OrderStatusConfirmed, models.OrderStatusCancelled, models.OrderStatusPending,
		models.OrderStatusConfirmed, models.OrderStatusConfirmed, models.OrderStatusConfirmed,
	}
	orders := make([]models.Order, len(orderIDs))
	for index, orderID := range orderIDs {
		orders[index] = models.Order{
			ID:         orderID,
			CustomerID: customerIDs[index],
			VehicleID:  vehicleIDs[index],
			OrderDate:  date(2026, 9, 1+index),
			Status:     statuses[index],
		}
	}
	if err := db.Create(&orders).Error; err != nil {
		return fmt.Errorf("creating orders: %w", err)
	}

	paymentStatuses := []string{
		models.PaymentStatusPaid, models.PaymentStatusPaid, models.PaymentStatusPending,
		models.PaymentStatusFailed, models.PaymentStatusPaid, models.PaymentStatusPaid,
		models.PaymentStatusPending, models.PaymentStatusPaid, models.PaymentStatusPaid,
		models.PaymentStatusPaid, models.PaymentStatusFailed, models.PaymentStatusPending,
		models.PaymentStatusPaid, models.PaymentStatusPaid, models.PaymentStatusPaid,
	}
	prices := []float64{1450000, 3800000, 1050000, 1850000, 1600000, 750000, 1350000, 1550000, 1700000, 1950000, 750000, 1050000, 1350000, 1550000, 1600000}
	payments := make([]models.Payment, len(orderIDs))
	for index, orderID := range orderIDs {
		payments[index] = models.Payment{
			ID:      uint(index + 1),
			OrderID: orderID,
			Amount:  prices[index],
			Status:  paymentStatuses[index],
		}
		if paymentStatuses[index] == models.PaymentStatusPaid {
			paidAt := dateTime(2026, 9, 1+index, 10, 30)
			payments[index].PaidAt = &paidAt
		}
	}
	if err := db.Create(&payments).Error; err != nil {
		return fmt.Errorf("creating payments: %w", err)
	}

	deliveryStatuses := []string{
		models.DeliveryStatusNotScheduled, models.DeliveryStatusDelivered, models.DeliveryStatusScheduled,
		models.DeliveryStatusNotScheduled, models.DeliveryStatusOutForDelivery, models.DeliveryStatusScheduled,
		models.DeliveryStatusNotScheduled, models.DeliveryStatusScheduled, models.DeliveryStatusDelivered,
		models.DeliveryStatusScheduled, models.DeliveryStatusNotScheduled, models.DeliveryStatusNotScheduled,
		models.DeliveryStatusScheduled, models.DeliveryStatusScheduled, models.DeliveryStatusScheduled,
	}
	deliveries := make([]models.Delivery, len(orderIDs))
	for index, orderID := range orderIDs {
		deliveries[index] = models.Delivery{
			ID:      uint(index + 1),
			OrderID: orderID,
			Status:  deliveryStatuses[index],
		}
		if deliveryStatuses[index] == models.DeliveryStatusScheduled || deliveryStatuses[index] == models.DeliveryStatusOutForDelivery {
			scheduledDate := dateTime(2026, 9, 10+index, 9, 0)
			deliveries[index].ScheduledDate = &scheduledDate
		}
		if deliveryStatuses[index] == models.DeliveryStatusDelivered {
			deliveredAt := dateTime(2026, 9, 10+index, 15, 30)
			deliveries[index].DeliveredAt = &deliveredAt
		}
	}
	if err := db.Create(&deliveries).Error; err != nil {
		return fmt.Errorf("creating deliveries: %w", err)
	}

	return nil
}

func clearDatabase(db *gorm.DB) error {
	for _, model := range []interface{}{
		&models.Delivery{}, &models.Payment{}, &models.Order{},
		&models.Vehicle{}, &models.Customer{},
	} {
		if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(model).Error; err != nil {
			return fmt.Errorf("clearing %T: %w", model, err)
		}
	}
	return nil
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func dateTime(year int, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}
