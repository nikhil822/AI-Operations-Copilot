package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const openRouterURL = "https://openrouter.ai/api/v1/chat/completions"

type OpenRouterProvider struct {
	APIKey     string
	HTTPClient *http.Client
	Model      string
}

func NewOpenRouterProvider(
	apiKey string,
	model string,
) *OpenRouterProvider {
	return &OpenRouterProvider{
		APIKey: apiKey,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OpenRouterProvider) Chat(
	ctx context.Context,
	request ChatRequest,
) (*ChatResponse, error) {

	if request.Model == "" {
		request.Model = p.Model
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal OpenRouter request: %w",
			err,
		)
	}

	const maxAttempts = 3

	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {

		response, retry, err := p.doRequest(
			ctx,
			body,
		)

		if err == nil {
			return response, nil
		}

		lastErr = err

		if !retry || attempt == maxAttempts {
			break
		}

		backoff := time.Duration(1<<(attempt-1)) * time.Second

		select {
		case <-ctx.Done():
			return nil, ctx.Err()

		case <-time.After(backoff):
		}
	}

	return nil, lastErr
}

func (p *OpenRouterProvider) doRequest(
	ctx context.Context,
	body []byte,
) (*ChatResponse, bool, error) {

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		openRouterURL,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, false, fmt.Errorf(
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

	req.Header.Set(
		"X-Title",
		"Cars24 AI Operations Copilot",
	)

	resp, err := p.HTTPClient.Do(req)

	if err != nil {
		// Network errors can be transient.
		return nil, true, fmt.Errorf(
			"OpenRouter request failed: %w",
			err,
		)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, true, fmt.Errorf(
			"failed to read OpenRouter response: %w",
			err,
		)
	}

	// Retry common transient failures.
	if resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode >= 500 {

		return nil, true, fmt.Errorf(
			"OpenRouter temporarily unavailable (HTTP %d): %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, false, fmt.Errorf(
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
		return nil, false, fmt.Errorf(
			"failed to parse OpenRouter response: %w",
			err,
		)
	}

	if result.Error != nil {
		message := result.Error.Message

		// Some upstream/provider errors can arrive as HTTP 200
		// with an error object, so retry obvious temporary errors.
		retry := strings.Contains(
			strings.ToLower(message),
			"temporarily",
		) || strings.Contains(
			strings.ToLower(message),
			"overloaded",
		)

		return nil, retry, fmt.Errorf(
			"OpenRouter API error: %s",
			message,
		)
	}

	if len(result.Choices) == 0 {
		return nil, false, fmt.Errorf(
			"OpenRouter returned no choices",
		)
	}

	return &result, false, nil
}