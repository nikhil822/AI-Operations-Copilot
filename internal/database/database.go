package database

import (
	"fmt"
	"os"
	"path/filepath"

	"ai-copilot/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Connect(databasePath string) (*gorm.DB, error) {
	dir := filepath.Dir(databasePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, fmt.Errorf(
				"failed to create database directory %q: %w",
				dir,
				err,
			)
		}
	}
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
