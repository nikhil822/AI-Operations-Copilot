package database

import (
	"fmt"

	"ai-copilot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Customer{},
		&models.Vehicle{},
		&models.Order{},
		&models.Payment{},
		&models.Delivery{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}