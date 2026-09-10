package models

import "time"

const (
	DeliveryStatusNotScheduled  = "not_scheduled"
	DeliveryStatusScheduled     = "scheduled"
	DeliveryStatusOutForDelivery = "out_for_delivery"
	DeliveryStatusDelivered     = "delivered"
)

type Delivery struct {
	ID            uint       `gorm:"primaryKey"`
	OrderID       uint       `gorm:"not null;uniqueIndex"`
	Status        string     `gorm:"not null;index"`
	ScheduledDate *time.Time
	DeliveredAt   *time.Time

	Order Order `gorm:"foreignKey:OrderID"`
}