package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabasePath     string
	AnthropicAPIKey  string
	AnthropicModel   string
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
		model = "claude-3-5-haiku-latest"
	}

	return Config{
		Port:            port,
		DatabasePath:    databasePath,
		AnthropicAPIKey: os.Getenv("API_KEY"),
		AnthropicModel:  model,
	}
}