package main

import (
	"log"

	"ai-copilot/internal/config"
	"ai-copilot/internal/database"
	"ai-copilot/internal/seed"
)

// Standalone entrypoint: `go run ./seed`. Useful for local development
// or forcing a fresh reseed. The server also seeds itself automatically
// on startup if the database is empty (see cmd/server/main.go), so
// running this manually is optional in most deployments.
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	if err := seed.Run(db); err != nil {
		log.Fatal(err)
	}

	log.Println("seed completed successfully")
}
