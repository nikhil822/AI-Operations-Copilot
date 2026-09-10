package main

import (
	"log"

	"ai-copilot/internal/config"
	"ai-copilot/internal/database"
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

	log.Printf("database initialized successfully")
	log.Printf("server will run on port %s", cfg.Port)
}