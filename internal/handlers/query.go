package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"ai-copilot/internal/ai"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	Agent *ai.Agent
}

type QueryRequest struct {
	Query string `json:"query"`
}

type QueryResponse struct {
	Answer string `json:"answer"`
}

func NewQueryHandler(agent *ai.Agent) *QueryHandler {
	return &QueryHandler{
		Agent: agent,
	}
}

func (h *QueryHandler) Query(c *gin.Context) {
	var req QueryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "request body must contain a valid JSON object with a query field",
		})
		return
	}

	query := strings.TrimSpace(req.Query)
	log.Printf(
		"AI query received: %s",
		query,
	)

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query is required",
		})
		return
	}

	// Prevent unnecessarily large prompts.
	if len(query) > 2000 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query is too long; maximum length is 2000 characters",
		})
		return
	}

	// Don't allow an LLM request to run forever.
	ctx, cancel := context.WithTimeout(c.Request.Context(), 45*time.Second)
	defer cancel()

	answer, err := h.Agent.Query(ctx, query)
	if err != nil {
		c.Error(err)

		log.Printf(
			"AI query failed: error=%v",
			err,
		)

		c.JSON(http.StatusBadGateway, gin.H{
			"error": "AI service temporarily unavailable",
		})
		return
	}

	log.Printf("AI query completed successfully")

	c.JSON(http.StatusOK, QueryResponse{
		Answer: answer,
	})
}
