package main

import (
	"log"
	"net/http"
	"time"

	"ai-copilot/internal/ai"
	"ai-copilot/internal/config"
	"ai-copilot/internal/database"
	"ai-copilot/internal/handlers"
	"ai-copilot/internal/tools"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	if cfg.APIKey == "" {
		log.Fatal("OPENROUTER_API_KEY is not configured")
	}

	// Connect to database.
	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create/update database tables.
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	// Create database-backed tools.
	orderTool := &tools.OrderTool{
		DB: db,
	}

	paymentTool := &tools.PaymentTool{
		DB: db,
	}

	deliveryTool := &tools.DeliveryTool{
		DB: db,
	}

	summaryTool := &tools.SummaryTool{
		DB: db,
	}

	// Create tool registry.
	registry := tools.NewRegistry(
		orderTool,
		paymentTool,
		deliveryTool,
		summaryTool,
	)

	// Create LLM provider.
	provider := ai.NewOpenRouterProvider(
		cfg.APIKey,
	)

	// Create AI agent.
	agent := ai.NewAgent(
		provider,
		registry,
		cfg.Model,
	)

	// Create HTTP handlers.
	queryHandler := handlers.NewQueryHandler(agent)

	// Create Gin router.
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// AI query endpoint.
	router.POST("/query", queryHandler.Query)

	// HTTP server.
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("Cars24 AI Operations Copilot listening on port %s", cfg.Port)
	log.Printf("Using OpenRouter model: %s", cfg.Model)

	if err := server.ListenAndServe(); err != nil &&
		err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
