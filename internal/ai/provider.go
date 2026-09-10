package ai

import "context"

type LLMProvider interface {
	Chat(
		ctx context.Context,
		request ChatRequest,
	) (*ChatResponse, error)
}