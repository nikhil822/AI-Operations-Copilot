package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabasePath string
	APIKey       string
	Model        string
}

func Load() Config {
	// Ignore the error because .env is optional.
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "./data/cars24.db"
	}

	model := os.Getenv("MODEL")
	if model == "" {
		model = "openrouter/free"
	}

	return Config{
		Port:            port,
		DatabasePath:    databasePath,
		APIKey: os.Getenv("API_KEY"),
		Model:  model,
	}
}
