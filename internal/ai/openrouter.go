package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouterProvider struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewOpenRouterProvider(apiKey string) *OpenRouterProvider {
	return &OpenRouterProvider{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *OpenRouterProvider) Chat(
	ctx context.Context,
	request ChatRequest,
) (*ChatResponse, error) {

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal OpenRouter request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		openRouterURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create OpenRouter request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+p.APIKey,
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	// Optional OpenRouter metadata.
	req.Header.Set(
		"X-Title",
		"Cars24 AI Operations Copilot",
	)

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"OpenRouter request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read OpenRouter response: %w",
			err,
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"OpenRouter returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var result ChatResponse

	if err := json.Unmarshal(
		responseBody,
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to parse OpenRouter response: %w",
			err,
		)
	}

	if result.Error != nil {
		return nil, fmt.Errorf(
			"OpenRouter API error: %s",
			result.Error.Message,
		)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf(
			"OpenRouter returned no choices",
		)
	}

	return &result, nil
}