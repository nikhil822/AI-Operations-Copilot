package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"ai-copilot/internal/ai"
	"ai-copilot/internal/config"
	"ai-copilot/internal/database"
	"ai-copilot/internal/handlers"
	"ai-copilot/internal/models"
	"ai-copilot/internal/seed"
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

	// Auto-seed on first boot. Deployments without a persistent disk
	// (e.g. Render's free tier) get a fresh, empty SQLite file on
	// every deploy/restart and typically have no shell access to run
	// `go run ./seed` manually — so the server seeds itself whenever
	// the orders table is empty, rather than requiring a manual step.
	var orderCount int64
	if err := db.Model(&models.Order{}).Count(&orderCount).Error; err != nil {
		log.Fatal(err)
	}

	if orderCount == 0 {
		log.Println("database is empty — running seed automatically")

		if err := seed.Run(db); err != nil {
			log.Fatal(err)
		}

		log.Println("seed completed successfully")
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
		cfg.Model,
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
	router.Use(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = fmt.Sprintf(
				"%d",
				time.Now().UnixNano(),
			)
		}

		c.Set("request_id", requestID)

		c.Header(
			"X-Request-ID",
			requestID,
		)

		c.Next()
	})

	// CORS — allows the frontend to be served separately (e.g. a
	// different port during local development) from the API.
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check.
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// AI query endpoint.
	router.POST("/query", queryHandler.Query)

	// Frontend — served from the same binary so the whole app is a
	// single deployable artifact. index.html at "/", assets under
	// "/static/*" (referenced as /static/style.css, /static/app.js).
	router.StaticFile("/", "./web/index.html")
	router.Static("/static", "./web")

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
