package main

import (
	"fmt"
	"log"
	"time"

	"ai-copilot/internal/config"
	"ai-copilot/internal/database"
	"ai-copilot/internal/models"

	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	if err := seed(db); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed successfully")
}

func seed(db *gorm.DB) error {
	// Clear existing data so the seed is repeatable.
	if err := clearDatabase(db); err != nil {
		return err
	}

	customers := []models.Customer{
		{
			ID:    1,
			Name:  "Rahul Sharma",
			Phone: "+91-9876543210",
			Email: "rahul.sharma@example.com",
		},
		{
			ID:    2,
			Name:  "Priya Mehta",
			Phone: "+91-9876543211",
			Email: "priya.mehta@example.com",
		},
		{
			ID:    3,
			Name:  "Amit Verma",
			Phone: "+91-9876543212",
			Email: "amit.verma@example.com",
		},
		{
			ID:    4,
			Name:  "Sneha Kapoor",
			Phone: "+91-9876543213",
			Email: "sneha.kapoor@example.com",
		},
		{
			ID:    5,
			Name:  "Arjun Singh",
			Phone: "+91-9876543214",
			Email: "arjun.singh@example.com",
		},
		{
			ID:    6,
			Name:  "Neha Gupta",
			Phone: "+91-9876543215",
			Email: "neha.gupta@example.com",
		},
		{
			ID:    7,
			Name:  "Vikram Malhotra",
			Phone: "+91-9876543216",
			Email: "vikram.malhotra@example.com",
		},
		{
			ID:    8,
			Name:  "Ananya Rao",
			Phone: "+91-9876543217",
			Email: "ananya.rao@example.com",
		},
		{
			ID:    9,
			Name:  "Karan Joshi",
			Phone: "+91-9876543218",
			Email: "karan.joshi@example.com",
		},
		{
			ID:    10,
			Name:  "Riya Das",
			Phone: "+91-9876543219",
			Email: "riya.das@example.com",
		},
	}

	if err := db.Create(&customers).Error; err != nil {
		return fmt.Errorf("creating customers: %w", err)
	}

	vehicles := []models.Vehicle{
		{
			ID:    1,
			Make:  "Maruti Suzuki",
			Model: "Swift",
			Year:  2023,
			Price: 750000,
		},
		{
			ID:    2,
			Make:  "Hyundai",
			Model: "Creta",
			Year:  2024,
			Price: 1450000,
		},
		{
			ID:    3,
			Make:  "Tata",
			Model: "Nexon",
			Year:  2023,
			Price: 1050000,
		},
		{
			ID:    4,
			Make:  "Mahindra",
			Model: "XUV700",
			Year:  2024,
			Price: 1850000,
		},
		{
			ID:    5,
			Make:  "Toyota",
			Model: "Fortuner",
			Year:  2022,
			Price: 3800000,
		},
		{
			ID:    6,
			Make:  "Honda",
			Model: "City",
			Year:  2023,
			Price: 1350000,
		},
		{
			ID:    7,
			Make:  "Kia",
			Model: "Seltos",
			Year:  2024,
			Price: 1600000,
		},
		{
			ID:    8,
			Make:  "Volkswagen",
			Model: "Taigun",
			Year:  2023,
			Price: 1550000,
		},
		{
			ID:    9,
			Make:  "Skoda",
			Model: "Slavia",
			Year:  2024,
			Price: 1700000,
		},
		{
			ID:    10,
			Make:  "MG",
			Model: "Hector",
			Year:  2023,
			Price: 1950000,
		},
	}

	if err := db.Create(&vehicles).Error; err != nil {
		return fmt.Errorf("creating vehicles: %w", err)
	}

	orders := []models.Order{
		// IMPORTANT EDGE CASE:
		// Payment is paid but delivery has not been scheduled.
		{
			ID:         1289,
			CustomerID: 1,
			VehicleID: 2,
			OrderDate:  date(2026, 9, 1),
			Status:     models.OrderStatusConfirmed,
		},

		// Fully completed order.
		{
			ID:         2231,
			CustomerID: 2,
			VehicleID: 5,
			OrderDate:  date(2026, 8, 15),
			Status:     models.OrderStatusCompleted,
		},

		// Payment pending even though delivery is scheduled.
		{
			ID:         4521,
			CustomerID: 3,
			VehicleID: 3,
			OrderDate:  date(2026, 9, 3),
			Status:     models.OrderStatusConfirmed,
		},

		// Payment failed.
		{
			ID:         3345,
			CustomerID: 4,
			VehicleID: 4,
			OrderDate:  date(2026, 9, 2),
			Status:     models.OrderStatusPending,
		},

		// Vehicle is currently out for delivery.
		{
			ID:         5678,
			CustomerID: 5,
			VehicleID: 7,
			OrderDate:  date(2026, 8, 28),
			Status:     models.OrderStatusConfirmed,
		},

		// Normal scheduled delivery.
		{
			ID:         6789,
			CustomerID: 6,
			VehicleID: 1,
			OrderDate:  date(2026, 9, 4),
			Status:     models.OrderStatusConfirmed,
		},

		// Pending payment and no delivery.
		{
			ID:         7890,
			CustomerID: 7,
			VehicleID: 6,
			OrderDate:  date(2026, 9, 5),
			Status:     models.OrderStatusPending,
		},

		{
			ID:         8901,
			CustomerID: 8,
			VehicleID: 8,
			OrderDate:  date(2026, 8, 25),
			Status:     models.OrderStatusConfirmed,
		},

		{
			ID:         9012,
			CustomerID: 9,
			VehicleID: 9,
			OrderDate:  date(2026, 8, 20),
			Status:     models.OrderStatusCompleted,
		},

		{
			ID:         9123,
			CustomerID: 10,
			VehicleID: 10,
			OrderDate:  date(2026, 9, 1),
			Status:     models.OrderStatusConfirmed,
		},

		{
			ID:         9345,
			CustomerID: 1,
			VehicleID: 1,
			OrderDate:  date(2026, 8, 30),
			Status:     models.OrderStatusCancelled,
		},

		{
			ID:         9456,
			CustomerID: 2,
			VehicleID: 3,
			OrderDate:  date(2026, 9, 6),
			Status:     models.OrderStatusPending,
		},

		{
			ID:         9567,
			CustomerID: 4,
			VehicleID: 6,
			OrderDate:  date(2026, 8, 29),
			Status:     models.OrderStatusConfirmed,
		},

		{
			ID:         9678,
			CustomerID: 6,
			VehicleID: 8,
			OrderDate:  date(2026, 9, 2),
			Status:     models.OrderStatusConfirmed,
		},

		{
			ID:         9789,
			CustomerID: 8,
			VehicleID: 7,
			OrderDate:  date(2026, 9, 3),
			Status:     models.OrderStatusConfirmed,
		},
	}

	if err := db.Create(&orders).Error; err != nil {
		return fmt.Errorf("creating orders: %w", err)
	}

	// Payment data.
	paid1 := dateTime(2026, 9, 1, 10, 30)
	paid2 := dateTime(2026, 8, 15, 14, 20)
	paid3 := dateTime(2026, 8, 28, 11, 10)
	paid4 := dateTime(2026, 9, 4, 9, 15)
	paid5 := dateTime(2026, 8, 25, 16, 40)
	paid6 := dateTime(2026, 8, 20, 12, 0)
	paid7 := dateTime(2026, 9, 1, 13, 45)
	paid8 := dateTime(2026, 8, 29, 15, 30)
	paid9 := dateTime(2026, 9, 2, 10, 0)
	paid10 := dateTime(2026, 9, 3, 17, 20)

	payments := []models.Payment{
		// #1289: PAID
		{
			ID:      1,
			OrderID: 1289,
			Amount:  1450000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid1,
		},

		// #2231: PAID
		{
			ID:      2,
			OrderID: 2231,
			Amount:  3800000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid2,
		},

		// #4521: PENDING
		{
			ID:      3,
			OrderID: 4521,
			Amount:  1050000,
			Status:  models.PaymentStatusPending,
		},

		// #3345: FAILED
		{
			ID:      4,
			OrderID: 3345,
			Amount:  1850000,
			Status:  models.PaymentStatusFailed,
		},

		// #5678: PAID
		{
			ID:      5,
			OrderID: 5678,
			Amount:  1600000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid3,
		},

		// #6789: PAID
		{
			ID:      6,
			OrderID: 6789,
			Amount:  750000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid4,
		},

		// #7890: PENDING
		{
			ID:      7,
			OrderID: 7890,
			Amount:  1350000,
			Status:  models.PaymentStatusPending,
		},

		// #8901: PAID
		{
			ID:      8,
			OrderID: 8901,
			Amount:  1550000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid5,
		},

		// #9012: PAID
		{
			ID:      9,
			OrderID: 9012,
			Amount:  1700000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid6,
		},

		// #9123: PAID
		{
			ID:      10,
			OrderID: 9123,
			Amount:  1950000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid7,
		},

		// #9345: FAILED
		{
			ID:      11,
			OrderID: 9345,
			Amount:  750000,
			Status:  models.PaymentStatusFailed,
		},

		// #9456: PENDING
		{
			ID:      12,
			OrderID: 9456,
			Amount:  1050000,
			Status:  models.PaymentStatusPending,
		},

		// #9567: PAID
		{
			ID:      13,
			OrderID: 9567,
			Amount:  1350000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid8,
		},

		// #9678: PAID
		{
			ID:      14,
			OrderID: 9678,
			Amount:  1550000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid9,
		},

		// #9789: PAID
		{
			ID:      15,
			OrderID: 9789,
			Amount:  1600000,
			Status:  models.PaymentStatusPaid,
			PaidAt:  &paid10,
		},
	}

	if err := db.Create(&payments).Error; err != nil {
		return fmt.Errorf("creating payments: %w", err)
	}

	// Delivery data.
	schedule6789 := dateTime(2026, 9, 10, 9, 0)
	schedule5678 := dateTime(2026, 9, 8, 9, 0)
	delivered2231 := dateTime(2026, 8, 20, 15, 30)
	schedule8901 := dateTime(2026, 9, 9, 10, 0)
	delivered9012 := dateTime(2026, 8, 25, 14, 0)
	schedule9123 := dateTime(2026, 9, 12, 11, 0)
	schedule9567 := dateTime(2026, 9, 11, 10, 0)
	schedule9678 := dateTime(2026, 9, 13, 9, 30)
	schedule9789 := dateTime(2026, 9, 14, 12, 0)

	deliveries := []models.Delivery{
		// #1289: KEY EDGE CASE
		// Paid, but delivery not scheduled.
		{
			ID:      1,
			OrderID: 1289,
			Status:  models.DeliveryStatusNotScheduled,
		},

		// #2231: DELIVERED
		{
			ID:          2,
			OrderID:     2231,
			Status:      models.DeliveryStatusDelivered,
			DeliveredAt: &delivered2231,
		},

		// #4521: SCHEDULED despite payment pending.
		{
			ID:            3,
			OrderID:       4521,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule6789,
		},

		// #3345: NOT SCHEDULED because payment failed.
		{
			ID:      4,
			OrderID: 3345,
			Status:  models.DeliveryStatusNotScheduled,
		},

		// #5678: OUT FOR DELIVERY
		{
			ID:            5,
			OrderID:       5678,
			Status:        models.DeliveryStatusOutForDelivery,
			ScheduledDate: &schedule5678,
		},

		// #6789: SCHEDULED
		{
			ID:            6,
			OrderID:       6789,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule6789,
		},

		// #7890: NOT SCHEDULED
		{
			ID:      7,
			OrderID: 7890,
			Status:  models.DeliveryStatusNotScheduled,
		},

		// #8901: SCHEDULED
		{
			ID:            8,
			OrderID:       8901,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule8901,
		},

		// #9012: DELIVERED
		{
			ID:          9,
			OrderID:     9012,
			Status:      models.DeliveryStatusDelivered,
			DeliveredAt: &delivered9012,
		},

		// #9123: SCHEDULED
		{
			ID:            10,
			OrderID:       9123,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule9123,
		},

		// #9345: Cancelled order, therefore no delivery.
		{
			ID:      11,
			OrderID: 9345,
			Status:  models.DeliveryStatusNotScheduled,
		},

		// #9456: Not scheduled.
		{
			ID:      12,
			OrderID: 9456,
			Status:  models.DeliveryStatusNotScheduled,
		},

		// #9567: Scheduled.
		{
			ID:            13,
			OrderID:       9567,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule9567,
		},

		// #9678: Scheduled.
		{
			ID:            14,
			OrderID:       9678,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule9678,
		},

		// #9789: Scheduled.
		{
			ID:            15,
			OrderID:       9789,
			Status:        models.DeliveryStatusScheduled,
			ScheduledDate: &schedule9789,
		},
	}

	if err := db.Create(&deliveries).Error; err != nil {
		return fmt.Errorf("creating deliveries: %w", err)
	}

	return nil
}

func clearDatabase(db *gorm.DB) error {
	// Delete children first because of relationships.
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Delivery{}).Error; err != nil {
		return err
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Payment{}).Error; err != nil {
		return err
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Order{}).Error; err != nil {
		return err
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Vehicle{}).Error; err != nil {
		return err
	}

	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(&models.Customer{}).Error; err != nil {
		return err
	}

	return nil
}

func date(year int, month int, day int) time.Time {
	return time.Date(
		year,
		time.Month(month),
		day,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func dateTime(year int, month int, day int, hour int, minute int) time.Time {
	return time.Date(
		year,
		time.Month(month),
		day,
		hour,
		minute,
		0,
		0,
		time.UTC,
	)
}