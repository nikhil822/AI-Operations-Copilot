package models

type Customer struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"not null"`
	Phone string
	Email string
	Orders []Order `gorm:"foreignKey:CustomerID"`
}