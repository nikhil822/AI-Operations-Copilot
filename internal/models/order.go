package models

import "time"

const (
	OrderStatusPending   = "pending"
	OrderStatusConfirmed = "confirmed"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

type Order struct {
	ID         uint      `gorm:"primaryKey"`
	CustomerID uint      `gorm:"not null;index"`
	VehicleID  uint      `gorm:"not null;index"`
	OrderDate  time.Time `gorm:"not null"`
	Status     string    `gorm:"not null;index"`

	Customer Customer `gorm:"foreignKey:CustomerID"`
	Vehicle  Vehicle  `gorm:"foreignKey:VehicleID"`

	Payment  *Payment  `gorm:"foreignKey:OrderID"`
	Delivery *Delivery `gorm:"foreignKey:OrderID"`
}