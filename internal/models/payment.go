package models

import "time"

const (
	PaymentStatusPending = "pending"
	PaymentStatusPaid    = "paid"
	PaymentStatusFailed  = "failed"
)

type Payment struct {
	ID      uint      `gorm:"primaryKey"`
	OrderID uint      `gorm:"not null;uniqueIndex"`
	Amount  float64   `gorm:"not null"`
	Status  string    `gorm:"not null;index"`
	PaidAt  *time.Time

	Order Order `gorm:"foreignKey:OrderID"`
}