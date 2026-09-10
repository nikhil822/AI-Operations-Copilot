package models

type Vehicle struct {
	ID    uint    `gorm:"primaryKey"`
	Make  string  `gorm:"not null"`
	Model string  `gorm:"not null"`
	Year  int     `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Orders []Order `gorm:"foreignKey:VehicleID"`
}